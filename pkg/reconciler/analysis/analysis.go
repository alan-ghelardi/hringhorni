package analysis

import (
	"context"
	"time"

	k8sclock "k8s.io/utils/clock"

	"github.com/hako/durafmt"
	"github.com/nubank/hringhorni/pkg/reconciler/events"

	rolloutsv1alpha1 "github.com/nubank/hringhorni/pkg/apis/rollouts/v1alpha1"
	analysisreconciler "github.com/nubank/hringhorni/pkg/client/injection/reconciler/rollouts/v1alpha1/analysis"
	rolloutslisters "github.com/nubank/hringhorni/pkg/client/listers/rollouts/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"
)

const (
	defaultRequeueDuration = 5 * time.Second
)

var (
	// clock allows us to control the clock for testing purposes.
	clock k8sclock.PassiveClock = k8sclock.RealClock{}
)

// Reconciler implements analysisreconciler.Interface for
// Analysis resources.
type Reconciler struct {

	// rolloutLister allows us to read Rollout objects from the indexer.
	rolloutLister rolloutslisters.RolloutLister

	analyzers []Analyzer
}

type Analyzer interface {
	Analyze(ctx context.Context, analysis *rolloutsv1alpha1.Analysis, rollout *rolloutsv1alpha1.Rollout) error
}

// Check that our Reconciler implements the required interfaces
var (
	_ analysisreconciler.Interface = (*Reconciler)(nil)
)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, analysis *rolloutsv1alpha1.Analysis) reconciler.Event {
	logger := logging.FromContext(ctx)

	logger.Info("Reconciling Analysis")

	if analysis.IsDone() {
		logger.Info("Analysis is done - no further actions to be performed")
		return controller.NewSkipKey(analysis.Namespace + "/" + analysis.Name)
	}

	previousCondition := analysis.Status.GetCondition(rolloutsv1alpha1.AnalysisConditionSucceeded)

	err := r.reconcileAnalysis(ctx, analysis)
	if err != nil {
		if isRequeueError, _ := controller.IsRequeueKey(err); !isRequeueError {
			analysis.Status.MarkFailed(rolloutsv1alpha1.AnalysisConditionSucceeded, rolloutsv1alpha1.AnalysisErrorReason, err.Error())
		}
	}

	events.Emit(ctx, previousCondition, analysis.Status.GetCondition(rolloutsv1alpha1.AnalysisConditionSucceeded), analysis)

	return err
}

func (r *Reconciler) reconcileAnalysis(ctx context.Context, analysis *rolloutsv1alpha1.Analysis) error {
	rollout, err := r.rolloutLister.Rollouts(analysis.Namespace).Get(analysis.Spec.RolloutRef)
	if err != nil {
		return controller.NewPermanentError(err)
	}

	if analysis.Status.StartedAt == nil {
		analysis.Status.StartedAt = &metav1.Time{Time: clock.Now()}
	}

	if isInProgress(analysis) || analysis.HasIntervalElapsed(clock) {
		for _, analyzer := range r.analyzers {
			if err := analyzer.Analyze(ctx, analysis, rollout); err != nil {
				return err
			}
		}
	}

	if failed, message := hasFailed(analysis); failed {
		analysis.Status.LastEvaluatedAt = &metav1.Time{Time: clock.Now()}
		analysis.Status.CompletedAt = analysis.Status.LastEvaluatedAt
		analysis.Status.MarkFailed(rolloutsv1alpha1.AnalysisConditionSucceeded, rolloutsv1alpha1.AnalysisFailedReason, "The analysis detected a failure. Cause: %s", message)
		return nil
	}

	if isInProgress(analysis) {
		analysis.Status.MarkUnknown(rolloutsv1alpha1.AnalysisConditionSucceeded, rolloutsv1alpha1.AnalysisInProgressReason, "Waiting for all analyses in the suite to complete")
		return controller.NewRequeueAfter(defaultRequeueDuration)
	}

	analysis.Status.LastEvaluatedAt = &metav1.Time{Time: clock.Now()}

	if !analysis.HasDurationElapsed(clock) {
		duration := analysis.Status.LastEvaluatedAt.Add(analysis.Spec.Interval.Duration).Sub(clock.Now())
		formattedDuration := durafmt.Parse(duration).LimitFirstN(2)
		if len(analysis.Status.Conditions) < 2 {
			analysis.Status.MarkUnknown(rolloutsv1alpha1.AnalysisConditionSucceeded, rolloutsv1alpha1.AnalysisInProgressReason, "Waiting %s to start analyzing the rollout", formattedDuration)
		} else {
			analysis.Status.MarkUnknown(rolloutsv1alpha1.AnalysisConditionSucceeded, rolloutsv1alpha1.AnalysisInProgressReason, "Waiting %s to run a new analyses suite", formattedDuration)
		}
		return controller.NewRequeueAfter(duration)
	}

	analysis.Status.CompletedAt = analysis.Status.LastEvaluatedAt
	analysis.Status.MarkSucceeded(rolloutsv1alpha1.AnalysisConditionSucceeded, rolloutsv1alpha1.AnalysisSucceededReason, "Analysis finished with success: no deviations or anomalies were detected")

	return nil
}

func hasFailed(analysis *rolloutsv1alpha1.Analysis) (bool, string) {
	for _, condition := range analysis.Status.GetConditions() {
		if condition.Type != rolloutsv1alpha1.AnalysisConditionSucceeded && condition.IsFalse() {
			return true, condition.GetMessage()
		}
	}
	return false, ""
}

func isInProgress(analysis *rolloutsv1alpha1.Analysis) bool {
	for _, condition := range analysis.Status.GetConditions() {
		if condition.Type != rolloutsv1alpha1.AnalysisConditionSucceeded && condition.IsUnknown() {
			return true
		}
	}
	return false
}
