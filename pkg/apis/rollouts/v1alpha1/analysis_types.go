package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/clock"

	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Analysis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AnalysisSpec   `json:"spec,omitempty"`
	Status            AnalysisStatus `json:"status,omitempty"`
}

var (
	// Check that Analysis can be validated and defaulted.
	_ kmeta.OwnerRefable = (*Analysis)(nil)
	// Check that the type conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*Analysis)(nil)
)

// AnalysisSpec describes the desired state of this Analysis object.
type AnalysisSpec struct {

	// Name of the application this analysis refers to.
	AppName string `json:"appName"`

	// How long the analysis will be performed during the rollout.
	Duration *metav1.Duration `json:"duration,omitempty"`

	Interval *metav1.Duration `json:"interval"`

	// Name of the Rollout under analysis.
	RolloutRef string `json:"rolloutRef"`

	// +optional
	Canary *CanaryAnalysis `json:"canary,omitempty"`
}

type CanaryAnalysis struct {
	MetricGroups []MetricGroup `json:"metricGroup"`
}

// MetricGroup groups related metrics.
type MetricGroup struct {

	// Human readable name of the group.
	Name string `json:"name"`

	// Weight of the group in the experiment's analysis.
	Weight *float64 `json:"weight,omitempty"`

	// Metrics to be evaluated during the experiment. At least one metric is
	// required.
	Metrics []Metric `json:"metrics"`
}

type Metric struct {

	// Human readable name of this metric.
	Name string `json:"name"`

	// Query to retrieve the actual metric ffrom the metrics provider.
	Query string `json:"query"`

	// Whether or not this metric is critical. If set to true and the metric
	// ffails during the analysis, the entire group fails even though other
	// metrics pass.
	// +optional
	Critical *bool `json:"critical,omitempty"`

	// How deviations between the control and experiment groups are analyzed
	// for this metric. Will the experiment fails when results from the
	// experiment group are higher, lower or differ in any direction from
	// the results obtained from the control group?
	// +optional
	FailOn AnalysisDirection `json:"failOn,omitempty"`
}

type AnalysisDirection string

// String implements fmt.Stringer.
func (a AnalysisDirection) String() string {
	return string(a)
}

const (
	DirectionIncrease AnalysisDirection = "Increase"
	DirectionDecrease AnalysisDirection = "Decrease"
	DirectionEither   AnalysisDirection = "Either"
)

// AnalysisStatus describes the current state of this Analysis.
type AnalysisStatus struct {
	duckv1.Status `json:",inline"`

	// Timestamp at which the analysis has started.
	// +optional
	StartedAt *metav1.Time `json:"startedAt,omitempty"`

	// Timestamp at which the analysis has been evaluated for the last time.
	// +optional
	LastEvaluatedAt *metav1.Time `json:"lastEvaluatedAt,omitempty"`

	// Timestamp at which the analysis has completed.
	// +optional
	CompletedAt *metav1.Time `json:"completedAt,omitempty"`

	// Status of the canary analysis.
	// +optional
	Canary *CanaryAnalysisStatus `json:"canary,omitempty"`
}

type CanaryAnalysisStatus struct {

	// Identifier of the experiment on the external platform responsible
	// for running the automated canary analysis (e.g. Kayenta).
	ExternalID *string `json:"externalID,omitempty"`

	RequestedAt *metav1.Time `json:"requestedAt,omitempty"`

	Results *runtime.RawExtension `json:"results,omitempty"`
}

func (c *CanaryAnalysisStatus) HasElapsed(interval *metav1.Duration, clock clock.PassiveClock) bool {
	if c.RequestedAt == nil {
		return false
	}
	return clock.Now().Sub(c.RequestedAt.Time) >= interval.Duration
}

// GetStatus retrieves the status of the resource. Implements the KRShaped interface.
func (a *Analysis) GetStatus() *duckv1.Status {
	return &a.Status.Status
}

func (a *Analysis) IsDone() bool {
	return !a.Status.GetCondition(AnalysisConditionSucceeded).IsUnknown()
}

// HasDurationElapsed returns true if the configured duration of this Analysis
// has elapsed.
func (a *Analysis) HasDurationElapsed(clock clock.PassiveClock) bool {
	if a.Status.StartedAt == nil {
		return false
	}
	return clock.Now().Sub(a.Status.StartedAt.Time) >= a.Spec.Duration.Duration
}

// HasIntervalElapsed returns true if the configured interval of this Analysis
// has elapsed since the last evaluation.
func (a *Analysis) HasIntervalElapsed(clock clock.PassiveClock) bool {
	if a.Status.StartedAt == nil {
		return false
	}
	startTime := a.Status.StartedAt
	if a.Status.LastEvaluatedAt != nil {
		startTime = a.Status.LastEvaluatedAt
	}
	return clock.Now().Sub(startTime.Time) >= a.Spec.Interval.Duration
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AnalysisList is a list of Analysis objects.
type AnalysisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Analysis `json:"items"`
}
