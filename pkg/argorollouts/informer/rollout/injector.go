package informer

import (
	"context"

	argorolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/client/informers/externalversions/rollouts/v1alpha1"
	"github.com/nubank/hringhorni/pkg/argorollouts/informerfactory"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformer(func(ctx context.Context) (context.Context, controller.Informer) {
		factory := informerfactory.Get(ctx)
		informer := factory.Argoproj().V1alpha1().Rollouts()
		return context.WithValue(ctx, informerKey{}, informer), informer.Informer()
	})
}

type informerKey struct {
}

func Get(ctx context.Context) argorolloutsv1alpha1.RolloutInformer {
	informer, ok := ctx.Value(informerKey{}).(argorolloutsv1alpha1.RolloutInformer)
	if !ok {
		logger := logging.FromContext(ctx)
		logger.Panic("Unable to fetch a Argo Rollouts informer from the provided Context")
	}
	return informer
}
