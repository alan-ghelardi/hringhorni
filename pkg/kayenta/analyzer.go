package kayenta

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	k8sclock "k8s.io/utils/clock"

	"github.com/go-openapi/strfmt"
	rolloutsv1alpha1 "github.com/nubank/hringhorni/pkg/apis/rollouts/v1alpha1"
	"github.com/nubank/hringhorni/pkg/kayenta/apis/models"
	kayentaclient "github.com/nubank/hringhorni/pkg/kayenta/client"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
)

var (
	// clock allows us to control the clock for testing purposes.
	clock k8sclock.PassiveClock = k8sclock.RealClock{}
)

type Analyzer struct {
	client *kayentaclient.KayentaClient
}

type analysisConfigurations map[string]map[string]any

type canary map[string]any

func (a *Analyzer) Analyze(ctx context.Context, analysis *rolloutsv1alpha1.Analysis, rollout *rolloutsv1alpha1.Rollout) error {
	logger := logging.FromContext(ctx)

	if analysis.Status.Canary == nil {
		analysis.Status.Canary = &rolloutsv1alpha1.CanaryAnalysisStatus{}
	}
	if analysis.Status.Canary.ExternalID == nil || analysis.Status.Canary.HasElapsed(analysis.Spec.Interval, clock) {
		logger.Info("Creating a new canary analysis on Kayenta")
		request := buildCanaryAdhocExecutionRequest(analysis, rollout)
		canaryExecutionResponse, err := a.client.CreateCanaryAnalysis(ctx, request)
		if err != nil {
			return err
		}
		canaryID := canaryExecutionResponse.CanaryExecutionID
		logger.Infow("Created a new canary analysis on Kayenta", zap.String("kayenta/canary-analysis-id", canaryID))
		analysis.Status.Canary.ExternalID = pointer.String(canaryID)
		analysis.Status.Canary.RequestedAt = &metav1.Time{Time: clock.Now()}
	}

	canaryExecutionStatusResponse, err := a.client.GetCanaryAnalysis(ctx, *analysis.Status.Canary.ExternalID)
	if err != nil {
		return err
	}

	if !canaryExecutionStatusResponse.Complete {
		analysis.Status.MarkUnknown(rolloutsv1alpha1.CanaryConditionSucceeded, rolloutsv1alpha1.AnalysisInProgressReason, "Waiting for canary's analysis results to be available on Kayenta")
		return nil
	}

	classification := canaryExecutionStatusResponse.Result.JudgeResult.Score.Classification
	classificationReason := canaryExecutionStatusResponse.Result.JudgeResult.Score.ClassificationReason
	message := fmt.Sprintf("Canary analysis classified as %s: %s. For further details, refer to status.canary.results", classification, classificationReason)
	switch strings.ToLower(classification) {
	case "error":
		analysis.Status.MarkFailed(rolloutsv1alpha1.CanaryConditionSucceeded, rolloutsv1alpha1.AnalysisErrorReason, message)

	case "fail":
		analysis.Status.MarkFailed(rolloutsv1alpha1.CanaryConditionSucceeded, rolloutsv1alpha1.AnalysisErrorReason, message)

	case "pass":
		analysis.Status.MarkSucceeded(rolloutsv1alpha1.CanaryConditionSucceeded, rolloutsv1alpha1.AnalysisSucceededReason, message)

	case "marginal", "nodata":
		analysis.Status.MarkFailed(rolloutsv1alpha1.CanaryConditionSucceeded, rolloutsv1alpha1.AnalysisInconclusiveReason, message)

	default:
		analysis.Status.MarkFailed(rolloutsv1alpha1.CanaryConditionSucceeded, rolloutsv1alpha1.AnalysisErrorReason, "Kayenta returned an unexpected classification for the canary analysis: %s", classification)
	}

	if canaryExecutionStatusResponse.Result != nil && canaryExecutionStatusResponse.Result.JudgeResult != nil {
		results, err := json.Marshal(canaryExecutionStatusResponse.Result.JudgeResult)
		if err != nil {
			return err
		}
		analysis.Status.Canary.Results = &runtime.RawExtension{Raw: results}
	}

	return nil
}

func buildCanaryAdhocExecutionRequest(analysis *rolloutsv1alpha1.Analysis, rollout *rolloutsv1alpha1.Rollout) *models.CanaryAdhocExecutionRequest {
	startTime := analysis.Status.StartedAt.UTC()
	endTime := startTime.Add(analysis.Spec.Duration.Duration)

	request := &models.CanaryAdhocExecutionRequest{
		CanaryConfig: &models.CanaryConfig{
			Name:         analysis.Spec.RolloutRef,
			Description:  fmt.Sprintf("canary configuration for the Rollout %s", analysis.Spec.RolloutRef),
			Applications: []string{analysis.Spec.AppName},
			// This is the only judge supported currently.
			Judge: &models.CanaryJudgeConfig{
				JudgeConfigurations: map[string]string{},
				Name:                pointer.String("NetflixACAJudge-v1.0"),
			},
			Metrics: []*models.CanaryMetricConfig{},
			Classifier: &models.CanaryClassifierConfig{
				GroupWeights: map[string]float64{},
			},
		},
		ExecutionRequest: &models.CanaryExecutionRequest{
			Scopes: map[string]models.CanaryScopePair{
				"default": {
					ControlScope: &models.CanaryScope{
						Location: "baseline",
						Scope:    *rollout.Status.Canary.BaselineDeployment,
						Start:    strfmt.DateTime(startTime),
						End:      strfmt.DateTime(endTime),
						Step:     1,
					},
					ExperimentScope: &models.CanaryScope{
						Location: "canary",
						Scope:    *rollout.Status.Canary.CanaryDeployment,
						Start:    strfmt.DateTime(startTime),
						End:      strfmt.DateTime(endTime),
						Step:     1,
					},
				},
			},
			// TODO(alan-ghelardi): make this configurable
			Thresholds: &models.CanaryClassifierThresholdsConfig{
				Marginal: 75.0,
				Pass:     95.0,
			},
		},
	}

	for _, group := range analysis.Spec.Canary.MetricGroups {
		request.CanaryConfig.Classifier.GroupWeights[group.Name] = *group.Weight
		for _, metric := range group.Metrics {
			canaryMetricConfig := &models.CanaryMetricConfig{
				Name: metric.Name,
				// Note: only Prometheus is supported nowadays.
				Query: &models.CanaryMetricSetQueryConfig{
					CustomInlineTemplate: fmt.Sprintf("PromQL:%s", relabelQuery(metric.Query)),
					ServiceType:          "prometheus",
					Type:                 "prometheus",
				},
				Groups: []string{group.Name},
				AnalysisConfigurations: analysisConfigurations{
					"canary": canary{
						"direction": strings.ToLower(metric.FailOn.String()),
						// TODO(alan-ghelardi): shall we make this configurable?
						"nanStrategy": "replace",
						"critical":    *metric.Critical,
					},
				},
				// This is the only scope supported by Kayenta.
				ScopeName: "default",
			}
			request.CanaryConfig.Metrics = append(request.CanaryConfig.Metrics, canaryMetricConfig)
		}
	}

	return request
}

// relabelQuery adds the location and scope labels to the promql query in
// question.
func relabelQuery(query string) string {
	return fmt.Sprintf(`label_replace(label_replace(%s, "location", "${location}", "", ""), "scope", "${scope}", "", "")`, query)
}

func New() *Analyzer {
	return &Analyzer{client: kayentaclient.New()}
}

type kayentaAnalyzerKey struct {
}

// Get returns the KayentaAnalyzer object attached to the provided Context
// or panics if none is present.
func Get(ctx context.Context) *Analyzer {
	analyzer, ok := ctx.Value(kayentaAnalyzerKey{}).(*Analyzer)
	if !ok {
		logging.FromContext(ctx).Panic("Cannot find a KayentaAnalyzer object in the provided Context")
	}
	return analyzer
}

// With attaches the KayentaAnalyzer object to the provided Context.
func With(ctx context.Context, analyzer *Analyzer) context.Context {
	return context.WithValue(ctx, kayentaAnalyzerKey{}, analyzer)
}

func init() {
	injection.Default.RegisterClient(func(ctx context.Context, _ *rest.Config) context.Context {
		return With(ctx, New())
	})
}
