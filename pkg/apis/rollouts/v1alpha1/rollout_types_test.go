package v1alpha1

import (
	"testing"
)

func TestIsRollingBack(t *testing.T) {
	assertIsRollingBack := func(t *testing.T, rollout *Rollout) {
		t.Helper()
		if !rollout.IsRollingBack() {
			t.Error("IsRollingBack: want true, but got false")
		}
	}

	assertIsNotRollingBack := func(t *testing.T, rollout *Rollout) {
		t.Helper()
		if rollout.IsRollingBack() {
			t.Error("IsRollingBack: want false, but got true")
		}
	}

	rollout := &Rollout{}
	assertIsNotRollingBack(t, rollout)

	rollout.Status.InitializeConditions()
	assertIsNotRollingBack(t, rollout)

	rollout.Status.MarkRollbackUnknown(RollbackInProgressReason, "lorem ipsum")
	assertIsRollingBack(t, rollout)

	rollout.Status.MarkRollbackSucceeded(RollbackSucceededReason, "lorem ipsum")
	assertIsNotRollingBack(t, rollout)
}
