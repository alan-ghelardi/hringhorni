package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
)

const (
	// RolloutConditionSucceeded identifies if the rollout has
	// finished.
	RolloutConditionSucceeded apis.ConditionType = "RolloutSucceeded"

	RollbackConditionSucceeded apis.ConditionType = "RollbackSucceeded"
)

// RolloutReason is the reason for each transition in the Rollout conditions
// throughout its lifecycle.
type RolloutReason string

// String implements Stringer.
func (r RolloutReason) String() string {
	return string(r)
}

const (

	// Reason set when the Rollout was created on the cluster but it hasn't
	// started yet.
	RolloutPendingReason RolloutReason = "RolloutPending"

	// Reason set when the rollout has started.
	RolloutInProgressReason RolloutReason = "RolloutInProgress"

	// Reason set when the rollout has completed with success.
	RolloutSucceededReason RolloutReason = "RolloutSucceeded"

	// Reason set when the rollout has completed with a failure.
	RolloutFailedReason RolloutReason = "RolloutFailed"

	// Reason set when the controller encounters an unrecoverable error
	// reconciling the Rollout.
	RolloutErrorReason RolloutReason = "RolloutError"

	// Reason set when the rollout has timed out.
	RolloutTimedOutReason RolloutReason = "RolloutTimedOut"

	// Reason set when the rollout is rolling back.
	RollbackInProgressReason RolloutReason = "RollbackInProgress"

	// Reason set when the rollback succeeds.
	RollbackSucceededReason RolloutReason = "RollbackSucceeded"

	// Reason set when the rollback failss.
	RollbackFailedReason RolloutReason = "RollbackFailed"

	// Reason set when the rollback has timed out.
	RollbackTimedOutReason RolloutReason = "RollbackTimedOut"

	// Reason set when the controller encounters an unrecoverable error
	// processing the rollback.
	RollbackErrorReason RolloutReason = "RollbackError"
)

var rolloutCondSet = apis.NewBatchConditionSet()

// GetGroupVersionKind implements kmeta.OwnerRefable
func (*Rollout) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("Rollout")
}

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (*Rollout) GetConditionSet() apis.ConditionSet {
	return rolloutCondSet
}

// InitializeConditions sets the initial values to the conditions.
func (r *RolloutStatus) InitializeConditions() {
	rolloutCondSet.Manage(r).InitializeConditions()
}

// MarkFailed signals an error in the state of the Rollout.
func (r *RolloutStatus) MarkFailed(reason RolloutReason, message string, args ...any) {
	rolloutCondSet.Manage(r).MarkFalse(
		RolloutConditionSucceeded,
		reason.String(),
		message,
		args...)
}

// MarkUnknown signals that the Rollout hasn't finished yet.
func (r *RolloutStatus) MarkUnknown(reason RolloutReason, message string, args ...any) {
	rolloutCondSet.Manage(r).MarkUnknown(
		RolloutConditionSucceeded,
		reason.String(),
		message,
		args...)
}

// MarkSucceeded signals that the Rollout has finished with success.
func (r *RolloutStatus) MarkSucceeded(reason RolloutReason, message string, args ...any) {
	rolloutCondSet.Manage(r).MarkTrueWithReason(RolloutConditionSucceeded, reason.String(), message, args...)
}

// MarkRollbackUnknown signals that the rollback started but hasn't finished
// yet.
func (r *RolloutStatus) MarkRollbackUnknown(reason RolloutReason, message string, args ...any) {
	rolloutCondSet.Manage(r).MarkUnknown(
		RollbackConditionSucceeded,
		reason.String(),
		message,
		args...)
}

// MarkRollbackSucceeded signals that the rollback has finished with success.
func (r *RolloutStatus) MarkRollbackSucceeded(reason RolloutReason, message string, args ...any) {
	rolloutCondSet.Manage(r).MarkTrueWithReason(RollbackConditionSucceeded, reason.String(), message, args...)
}

// MarkFailed signals an error in the rollback.
func (r *RolloutStatus) MarkRollbackFailed(reason RolloutReason, message string, args ...any) {
	rolloutCondSet.Manage(r).MarkFalse(
		RollbackConditionSucceeded,
		reason.String(),
		message,
		args...)
}
