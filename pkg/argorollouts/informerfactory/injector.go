package informerfactory

import (
	"context"
	"time"

	argorolloutsinformer "github.com/argoproj/argo-rollouts/pkg/client/informers/externalversions"
	"github.com/nubank/hringhorni/pkg/argorollouts/clientset"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformerFactory(func(ctx context.Context) context.Context {
		client := clientset.Get(ctx)
		factory := argorolloutsinformer.NewSharedInformerFactoryWithOptions(client, 1*time.Hour)
		return context.WithValue(ctx, informerFactoryKey{}, factory)
	})
}

type informerFactoryKey struct {
}

func Get(ctx context.Context) argorolloutsinformer.SharedInformerFactory {
	factory, ok := ctx.Value(informerFactoryKey{}).(argorolloutsinformer.SharedInformerFactory)
	if !ok {
		logger := logging.FromContext(ctx)
		logger.Panic("Unable to fetch a Argo Rollouts informer factory from the provided Context")
	}
	return factory
}
