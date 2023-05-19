package rollout

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	rolloutsv1alpha1 "github.com/nubank/hringhorni/pkg/apis/rollouts/v1alpha1"
	"github.com/nubank/hringhorni/pkg/deployments"
	"github.com/nubank/hringhorni/pkg/reconciler/events"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"

	"k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/polymorphichelpers"
)

func (r *Reconciler) rollBack(ctx context.Context, rollout *rolloutsv1alpha1.Rollout, productionDeployment *appsv1.Deployment) error {
	previousCondition := rollout.Status.GetCondition(rolloutsv1alpha1.RollbackConditionSucceeded)
	defer func() {
		events.Emit(ctx, previousCondition, rollout.Status.GetCondition(rolloutsv1alpha1.RollbackConditionSucceeded), rollout)
	}()

	if done, err := r.undoCanarySetup(ctx, rollout, productionDeployment); err != nil {
		return err
	} else if !done {
		// Wait for the next reconciliation cycle.
		return nil
	}

	if rollout.Status.RollBackRequested || rollout.Spec.Actions.RollBack {
		if err := r.RollBackToPreviousRevision(ctx, rollout, productionDeployment); err != nil {
			return err
		}
	}

	rollout.Status.MarkRollbackSucceeded(rolloutsv1alpha1.RollbackSucceededReason, "Rollout of revision %s has been successfully undone", rollout.Spec.Revision)

	return nil
}

func (r *Reconciler) undoCanarySetup(ctx context.Context, rollout *rolloutsv1alpha1.Rollout, productionDeployment *appsv1.Deployment) (bool, error) {
	logger := logging.FromContext(ctx)
	if productionDeployment != nil && rollout.Spec.Canary != nil {
		logger.Debug("Rolling back canary's setup")
		if done, err := r.canary.Undo(ctx, rollout, productionDeployment); err != nil {
			return done, err
		} else if !done {
			return done, nil
		}
	}
	return true, nil
}

func (r *Reconciler) RollBackToPreviousRevision(ctx context.Context, rollout *rolloutsv1alpha1.Rollout, productionDeployment *appsv1.Deployment) error {
	logger := logging.FromContext(ctx)

	if productionDeployment == nil {
		return controller.NewPermanentError(errors.New("unable to start rolling back the Deployment: no previous revision found"))
	}

	logger.Debugf("Rolling back Deployment %s to previous revision", productionDeployment.Name)

	if !rollout.Status.RollBackRequested {
		rollbacker, err := polymorphichelpers.RollbackerFor(schema.GroupKind{Group: appsv1.GroupName, Kind: "Deployment"}, r.kubeClient)
		if err != nil {
			return controller.NewPermanentError(err)
		}

		toRevision, err := strconv.ParseInt(rollout.Status.PreviousRollout.DeploymentRevision, 10, 0)
		if err != nil {
			return controller.NewPermanentError(fmt.Errorf("bad Deployment revision %q - the revision must be a valid integer", rollout.Status.PreviousRollout.Revision))
		}

		_, err = rollbacker.Rollback(productionDeployment, map[string]string{}, toRevision, util.DryRunNone)
		if err != nil {
			return controller.NewPermanentError(err)
		}
		rollout.Status.RollBackRequested = true
		rollout.Status.MarkRollbackUnknown(rolloutsv1alpha1.RollbackInProgressReason, "Rolling back to revision %s: requesting revision %d on Deployment %s",
			rollout.Status.PreviousRollout.Revision, toRevision, productionDeployment.Name)
		return controller.NewRequeueAfter(defaultRequeueDuration)
	}

	if ready, err := deployments.IsReady(productionDeployment, rollout.Status.PreviousRollout.Replicas); err != nil {
		return err
	} else if !ready {
		rollout.Status.MarkRollbackUnknown(rolloutsv1alpha1.RollbackInProgressReason, "Rolling back to revision %s: waiting for Deployment %s to meet the desired state",
			rollout.Status.PreviousRollout.Revision, productionDeployment.Name)
		return controller.NewRequeueAfter(defaultRequeueDuration)
	}

	return nil
}
