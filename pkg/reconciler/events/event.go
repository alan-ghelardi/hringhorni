package events

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

func Emit(ctx context.Context, previousCondition, newCondition *apis.Condition, obj runtime.Object) {
	logger := logging.FromContext(ctx)

	if newCondition == nil {
		logger.Debug("No event will be recorded because the new object's condition is nil")
		return
	}

	if previousCondition != nil &&
		previousCondition.Status == newCondition.Status &&
		previousCondition.GetReason() == newCondition.GetReason() &&
		previousCondition.GetMessage() == newCondition.GetMessage() {
		logger.Debug("No event will be recorded because there are no changes in the object's conditions")
		return
	}

	recorder := controller.GetEventRecorder(ctx)
	switch newCondition.Status {
	case corev1.ConditionTrue, corev1.ConditionUnknown:
		recorder.Event(obj, corev1.EventTypeNormal, newCondition.Reason, newCondition.Message)
	case corev1.ConditionFalse:
		recorder.Event(obj, corev1.EventTypeWarning, newCondition.Reason, newCondition.Message)
	}
	logger.Infof("Recorded event: type=%s, reason=%s, message=%s", newCondition.Type, newCondition.Reason, newCondition.Message)
}
