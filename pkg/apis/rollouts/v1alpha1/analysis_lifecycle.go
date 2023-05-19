package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
)

const (
	// AnalysisConditionSucceeded identifies if the analysis has finished.
	AnalysisConditionSucceeded apis.ConditionType = apis.ConditionSucceeded

	CanaryConditionSucceeded apis.ConditionType = "CanarySucceeded"
)

// AnalysisReason is the reason for each transition in the Analysis conditions
// throughout its lifecycle.
type AnalysisReason string

// String implements Stringer.
func (a AnalysisReason) String() string {
	return string(a)
}

const (

	// Reason set when the analysis has started, but the final result isn't
	// available yet.
	AnalysisInProgressReason AnalysisReason = "AnalysisInProgress"

	// Reason set when the analysis finished with a successful result.
	AnalysisSucceededReason AnalysisReason = "AnalysisSucceeded"

	// Reason set when the analysis has finished with a failure.
	AnalysisFailedReason AnalysisReason = "AnalysisFailed"

	// Reason set when the analysis produced inconclusive resultss.
	AnalysisInconclusiveReason AnalysisReason = "AnalysisInconclusive"

	// Reason set when the controller encounters an error reconciling the
	// Analysis object.
	AnalysisErrorReason AnalysisReason = "AnalysisError"
)

var analysisCondSet = apis.NewBatchConditionSet()

// GetGroupVersionKind implements kmeta.OwnerRefable
func (*Analysis) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("Analysis")
}

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (*Analysis) GetConditionSet() apis.ConditionSet {
	return analysisCondSet
}

// InitializeConditions sets the initial values to the conditions.
func (a *AnalysisStatus) InitializeConditions() {
	analysisCondSet.Manage(a).InitializeConditions()
}

// MarkFailed signals an error in the state of the Analysis.
func (a *AnalysisStatus) MarkFailed(condition apis.ConditionType, reason AnalysisReason, message string, args ...any) {
	analysisCondSet.Manage(a).MarkFalse(condition, reason.String(), message, args...)
}

// MarkUnknown signals that the Analysis hasn't finished yet.
func (a *AnalysisStatus) MarkUnknown(condition apis.ConditionType, reason AnalysisReason, message string, args ...any) {
	analysisCondSet.Manage(a).MarkUnknown(condition, reason.String(), message, args...)
}

// MarkSucceeded signals that the Analysis has finished with success.
func (a *AnalysisStatus) MarkSucceeded(condition apis.ConditionType, reason AnalysisReason, message string, args ...any) {
	analysisCondSet.Manage(a).MarkTrueWithReason(condition, reason.String(), message, args...)
}
