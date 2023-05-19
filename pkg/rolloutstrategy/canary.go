package rolloutstrategy

import (
	"context"
	"fmt"
	"math"
	"time"

	deploymentinformer "knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"
	"knative.dev/pkg/logging"

	appsv1 "k8s.io/api/apps/v1"

	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/utils/pointer"
	"knative.dev/pkg/controller"

	rolloutsv1alpha1 "github.com/nubank/hringhorni/pkg/apis/rollouts/v1alpha1"
	"github.com/nubank/hringhorni/pkg/deployments"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

const (
	defaultRequeueDuration = 5 * time.Second
)

type Canary struct {

	// deploymentLister allows us to read Deployments from the indexer.
	deploymentLister appslisters.DeploymentLister
}

func (c *Canary) ProgressRollout(ctx context.Context, rollout *rolloutsv1alpha1.Rollout, productionDeployment *appsv1.Deployment) (bool, error) {
	logger := logging.FromContext(ctx)

	if rollout.Status.Canary == nil {
		rollout.Status.Canary = &rolloutsv1alpha1.CanaryStatus{}
	}

	desiredReplicas := int32(math.Ceil(float64(*productionDeployment.Spec.Replicas) * (*rollout.Spec.Canary.Percentage / 100)))

	baseline, err := c.ensureBaselineDeployment(ctx, rollout, productionDeployment, desiredReplicas)
	if err != nil {
		return false, err
	}

	if ready, err := deployments.IsReady(baseline, desiredReplicas); err != nil {
		return false, err
	} else if !ready {
		logger.Debugf("Baseline Deployment %s isn't ready yet - waiting for the next reconciliation cycle", baseline.Name)
		rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutInProgressReason, "Rolling out revision %s (setting up canary): waiting for baseline Deployment (%s) to meet the desired conditions (%d of %d replicas are available)",
			rollout.Spec.Revision, baseline.Name, baseline.Status.AvailableReplicas, desiredReplicas)
		return false, nil
	}

	canary, err := c.ensureCanaryDeployment(ctx, rollout, desiredReplicas)
	if err != nil {
		return false, err
	}

	if ready, err := deployments.IsReady(canary, desiredReplicas); err != nil {
		return false, err
	} else if !ready {
		logger.Debugf("Canary Deployment %s isn't ready yet - waiting for the next reconciliation cycle", canary.Name)
		rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutInProgressReason, "Rolling out revision %s (setting up canary): waiting for canary Deployment (%s) to meet the desired conditions (%d of %d replicas are available)",
			rollout.Spec.Revision, canary.Name, canary.Status.AvailableReplicas, desiredReplicas)
		return false, nil
	}

	downscaleReplicas := rollout.Status.PreviousRollout.Replicas - desiredReplicas

	if downscaleReplicas != *productionDeployment.Spec.Replicas {
		logger.Debugf("Downscaling the production Deployment's replicas to %d", downscaleReplicas)
		if err := deployments.Scale(ctx, productionDeployment, downscaleReplicas); err != nil {
			return false, err
		}
	}

	if ready, err := deployments.IsReady(productionDeployment, downscaleReplicas); err != nil {
		return false, err
	} else if !ready {
		logger.Debugf("Production Deployment %s isn't ready yet - requeuing to process again", productionDeployment.Name)
		rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutInProgressReason, "Rolling out revision %s (setting up canary): waiting for production Deployment (%s) to meet the desired conditions (%d of %d replicas are available)",
			rollout.Spec.Revision, productionDeployment.Name, productionDeployment.Status.AvailableReplicas, downscaleReplicas)
		return false, controller.NewRequeueAfter(defaultRequeueDuration)
	}

	logger.Debug("Canary's setup is all set")

	return true, nil
}

func (c *Canary) ensureBaselineDeployment(ctx context.Context, rollout *rolloutsv1alpha1.Rollout, productionDeployment *appsv1.Deployment, desiredReplicas int32) (*appsv1.Deployment, error) {
	// Create a new baseline Deployment if it doesn't exist yet.
	if rollout.Status.Canary.BaselineDeployment == nil {
		baseline := deployments.NewBaseline(rollout, productionDeployment, desiredReplicas)
		createdBaseline, err := deployments.Create(ctx, baseline)
		if createdBaseline != nil {
			rollout.Status.Canary.BaselineDeployment = pointer.String(createdBaseline.Name)
		}
		if apierrors.IsAlreadyExists(err) {
			// Maybe the Rollout is somehow out of sync. Let's
			// requeue and try again as soon as possible.
			rollout.Status.Canary.BaselineDeployment = pointer.String(baseline.Name)
			rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutInProgressReason, "Rolling out revision %s (setting up canary): waiting for Deployment %s to be available in the indexer", rollout.Spec.Revision, baseline.Name)
			return nil, controller.NewRequeueAfter(defaultRequeueDuration)
		}
		return createdBaseline, err
	}

	// Most likely the baseline Deployment already exists. Let's read it
	// from the indexer.
	return c.getDeployment(rollout, *rollout.Status.Canary.BaselineDeployment)
}

func (c *Canary) getDeployment(rollout *rolloutsv1alpha1.Rollout, name string) (*appsv1.Deployment, error) {
	deployment, err := c.deploymentLister.Deployments(rollout.Namespace).Get(name)
	if err != nil {
		// Is the index somehow out of sync?
		if apierrors.IsNotFound(err) {
			rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutInProgressReason, "Rolling out revision %s (setting up canary): waiting for Deployment %s to be available in the indexer", rollout.Spec.Revision, name)
			return nil, controller.NewRequeueAfter(defaultRequeueDuration)
		}
		return nil, fmt.Errorf("error reading Deployment %s from the indexer: %w", name, err)
	}
	return deployment, nil
}

func (c *Canary) ensureCanaryDeployment(ctx context.Context, rollout *rolloutsv1alpha1.Rollout, desiredReplicas int32) (*appsv1.Deployment, error) {
	// Create a new canary Deployment if it doesn't exist yet.
	if rollout.Status.Canary.CanaryDeployment == nil {
		canary := deployments.NewCanary(rollout, desiredReplicas)
		createdCanary, err := deployments.Create(ctx, canary)
		if createdCanary != nil {
			rollout.Status.Canary.CanaryDeployment = pointer.String(createdCanary.Name)
		}
		if apierrors.IsAlreadyExists(err) {
			// Maybe the Rollout is somehow out of sync. Let's
			// requeue and try again as soon as possible.
			rollout.Status.Canary.CanaryDeployment = pointer.String(canary.Name)
			rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutInProgressReason, "Rolling out revision %s (setting up canary): waiting for Deployment %s to be available in the indexer", rollout.Spec.Revision, canary.Name)
			return nil, controller.NewRequeueAfter(defaultRequeueDuration)
		}
		return createdCanary, err
	}

	// Most likely the canary Deployment already exists. Let's read it
	// from the indexer.
	return c.getDeployment(rollout, *rollout.Status.Canary.CanaryDeployment)
}

func (c *Canary) Undo(ctx context.Context, rollout *rolloutsv1alpha1.Rollout, productionDeployment *appsv1.Deployment) (bool, error) {
	logger := logging.FromContext(ctx)

	if done, err := c.destroyDeployment(ctx, rollout.Namespace, *rollout.Status.Canary.CanaryDeployment); err != nil {
		return false, err
	} else if !done {
		markUnknown(rollout, "destroying canary Deployment %s", *rollout.Status.Canary.CanaryDeployment)
		return false, nil
	}

	if desiredReplicas := rollout.Status.PreviousRollout.Replicas; desiredReplicas > *productionDeployment.Spec.Replicas {
		logger.Infof("Upscaling production Deployment %s to match %d replicas", productionDeployment.Name, desiredReplicas)
		if err := deployments.Scale(ctx, productionDeployment, desiredReplicas); err != nil {
			return false, err
		}
	}

	if ready, err := deployments.IsReady(productionDeployment, rollout.Status.PreviousRollout.Replicas); err != nil {
		return false, err
	} else if !ready {
		logger.Debugf("Production Deployment %s isn't ready yet - requeuing to process again", productionDeployment.Name)
		markUnknown(rollout, "restoring Deployment %s to its original state", productionDeployment.Name)
		return false, controller.NewRequeueAfter(defaultRequeueDuration)
	}

	if done, err := c.destroyDeployment(ctx, rollout.Namespace, *rollout.Status.Canary.BaselineDeployment); err != nil {
		return false, err
	} else if !done {
		markUnknown(rollout, "destroying baseline Deployment %s", *rollout.Status.Canary.BaselineDeployment)
		return false, nil
	}

	return true, nil
}

func (c *Canary) destroyDeployment(ctx context.Context, namespace, name string) (bool, error) {
	logger := logging.FromContext(ctx)

	deployment, err := c.deploymentLister.Deployments(namespace).Get(name)
	if err != nil {
		// If the Deployment in question is no longer available we can interrupt the flow here.
		if apierrors.IsNotFound(err) {
			logger.Debugf("Deployment %s could not be found in the namespace %s - ignoring", namespace, name)
			return true, nil
		}
		return false, fmt.Errorf("error reading Deployment %s from the indexer: %w", name, err)
	}

	if *deployment.Spec.Replicas != 0 {
		logger.Debugf("Downscaling Deployment %s to 0 replicas", deployment.Name)
		if err := deployments.Scale(ctx, deployment, 0); err != nil {
			return false, err
		}
	}

	if ready, err := deployments.IsReady(deployment, 0); err != nil {
		return false, err
	} else if !ready {
		logger.Debugf("Deployment %s isn't ready yet - waiting for the next reconciliation cycle", deployment.Name)
		return false, nil
	}

	logger.Debugf("Deleting Deployment %s", deployment.Name)
	if err := deployments.Delete(ctx, deployment); err != nil {
		return false, err
	}

	return true, nil
}

func New(ctx context.Context) *Canary {
	return &Canary{
		deploymentLister: deploymentinformer.Get(ctx).Lister(),
	}
}

func markUnknown(rollout *rolloutsv1alpha1.Rollout, message string, args ...any) {
	args = append([]any{rollout.Spec.Revision}, args...)
	if rollout.IsRollingBack() {
		rollout.Status.MarkRollbackUnknown(rolloutsv1alpha1.RollbackInProgressReason, "Aborting rollout of revision %s (rolling back canary experiment): "+message, args...)
	} else {
		rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutInProgressReason, "Rolling out revision %s (reverting canary experiment): "+message, args...)
	}
}
