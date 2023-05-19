package events

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	rolloutsv1alpha1 "github.com/nubank/hringhorni/pkg/apis/rollouts/v1alpha1"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

func TestEmit(t *testing.T) {
	tests := []struct {
		name              string
		previousCondition *apis.Condition
		newCondition      *apis.Condition
		want              string
	}{{
		name: "emit event on True condition",
		newCondition: &apis.Condition{
			Status:  corev1.ConditionTrue,
			Reason:  "x",
			Message: "Lorem ipsum",
		},
		want: "Normal x Lorem ipsum",
	},
		{
			name: "emit event on False condition",
			newCondition: &apis.Condition{
				Status:  corev1.ConditionFalse,
				Reason:  "x",
				Message: "Lorem ipsum",
			},
			want: "Warning x Lorem ipsum",
		},
		{
			name: "emit event on Unknown condition",
			newCondition: &apis.Condition{
				Status:  corev1.ConditionUnknown,
				Reason:  "x",
				Message: "Lorem ipsum",
			},
			want: "Normal x Lorem ipsum",
		},
		{
			name: "do not emit the event if the conditions are semantically equal",
			previousCondition: &apis.Condition{
				Status:  corev1.ConditionTrue,
				Reason:  "x",
				Message: "Lorem ipsum",
			},
			newCondition: &apis.Condition{
				Status:  corev1.ConditionTrue,
				Reason:  "x",
				Message: "Lorem ipsum",
			},
		},
		{
			name: "emit the event if reazons and statuses are the same, but the message has changed",
			previousCondition: &apis.Condition{
				Status:  corev1.ConditionTrue,
				Reason:  "x",
				Message: "Lorem ipsum",
			},
			newCondition: &apis.Condition{
				Status:  corev1.ConditionTrue,
				Reason:  "x",
				Message: "Dolor sit amet",
			},
			want: "Normal x Dolor sit amet",
		},
		{
			name: "emit the event if the condition's status has changed",
			previousCondition: &apis.Condition{
				Status:  corev1.ConditionFalse,
				Reason:  "x",
				Message: "Lorem ipsum",
			},
			newCondition: &apis.Condition{
				Status:  corev1.ConditionTrue,
				Reason:  "x",
				Message: "Lorem ipsum",
			},
			want: "Normal x Lorem ipsum",
		},
		{
			name: "emit the event if the condition's reason has changed",
			previousCondition: &apis.Condition{
				Status:  corev1.ConditionFalse,
				Reason:  "x",
				Message: "Lorem ipsum",
			},
			newCondition: &apis.Condition{
				Status:  corev1.ConditionFalse,
				Reason:  "y",
				Message: "Lorem ipsum",
			},
			want: "Warning y Lorem ipsum",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eventRecorder := record.NewFakeRecorder(1)
			ctx := logging.WithLogger(context.Background(), zaptest.NewLogger(t).Sugar())
			ctx = controller.WithEventRecorder(ctx, eventRecorder)
			Emit(ctx, test.previousCondition, test.newCondition, &rolloutsv1alpha1.Rollout{})
			got := readEvent(eventRecorder)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func readEvent(eventRecorder *record.FakeRecorder) string {
	select {
	case event := <-eventRecorder.Events:
		return event

	default:
		return ""
	}
}
