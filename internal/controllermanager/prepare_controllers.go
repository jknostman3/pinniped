// Copyright 2020 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package controllermanager

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	k8sinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2/klogr"
	aggregatorclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"

	pinnipedv1alpha1 "github.com/suzerain-io/pinniped/generated/1.19/apis/pinniped/v1alpha1"
	pinnipedclientset "github.com/suzerain-io/pinniped/generated/1.19/client/clientset/versioned"
	pinnipedinformers "github.com/suzerain-io/pinniped/generated/1.19/client/informers/externalversions"
	"github.com/suzerain-io/pinniped/internal/controller/apicerts"
	"github.com/suzerain-io/pinniped/internal/controller/identityprovider/idpcache"
	"github.com/suzerain-io/pinniped/internal/controller/identityprovider/webhookcachecleaner"
	"github.com/suzerain-io/pinniped/internal/controller/identityprovider/webhookcachefiller"
	"github.com/suzerain-io/pinniped/internal/controller/issuerconfig"
	"github.com/suzerain-io/pinniped/internal/controllerlib"
	"github.com/suzerain-io/pinniped/internal/provider"
)

const (
	singletonWorker       = 1
	defaultResyncInterval = 3 * time.Minute
)

// Prepare the controllers and their informers and return a function that will start them when called.
func PrepareControllers(
	serverInstallationNamespace string,
	discoveryURLOverride *string,
	dynamicCertProvider provider.DynamicTLSServingCertProvider,
	servingCertDuration time.Duration,
	servingCertRenewBefore time.Duration,
	idpCache *idpcache.Cache,
) (func(ctx context.Context), error) {
	// Create k8s clients.
	k8sClient, aggregatorClient, pinnipedClient, err := createClients()
	if err != nil {
		return nil, fmt.Errorf("could not create clients for the controllers: %w", err)
	}

	// Create informers. Don't forget to make sure they get started in the function returned below.
	kubePublicNamespaceK8sInformers, installationNamespaceK8sInformers, installationNamespacePinnipedInformers :=
		createInformers(serverInstallationNamespace, k8sClient, pinnipedClient)

	// This string must match the name of the Service declared in the deployment yaml.
	const serviceName = "pinniped-api"

	// Create controller manager.
	controllerManager := controllerlib.
		NewManager().
		WithController(
			issuerconfig.NewPublisherController(
				serverInstallationNamespace,
				discoveryURLOverride,
				pinnipedClient,
				kubePublicNamespaceK8sInformers.Core().V1().ConfigMaps(),
				installationNamespacePinnipedInformers.Crd().V1alpha1().CredentialIssuerConfigs(),
				controllerlib.WithInformer,
			),
			singletonWorker,
		).
		WithController(
			apicerts.NewCertsManagerController(
				serverInstallationNamespace,
				k8sClient,
				installationNamespaceK8sInformers.Core().V1().Secrets(),
				controllerlib.WithInformer,
				controllerlib.WithInitialEvent,
				servingCertDuration,
				"Pinniped CA",
				serviceName,
			),
			singletonWorker,
		).
		WithController(
			apicerts.NewAPIServiceUpdaterController(
				serverInstallationNamespace,
				pinnipedv1alpha1.SchemeGroupVersion.Version+"."+pinnipedv1alpha1.GroupName,
				aggregatorClient,
				installationNamespaceK8sInformers.Core().V1().Secrets(),
				controllerlib.WithInformer,
			),
			singletonWorker,
		).
		WithController(
			apicerts.NewCertsObserverController(
				serverInstallationNamespace,
				dynamicCertProvider,
				installationNamespaceK8sInformers.Core().V1().Secrets(),
				controllerlib.WithInformer,
			),
			singletonWorker,
		).
		WithController(
			apicerts.NewCertsExpirerController(
				serverInstallationNamespace,
				k8sClient,
				installationNamespaceK8sInformers.Core().V1().Secrets(),
				controllerlib.WithInformer,
				servingCertRenewBefore,
			),
			singletonWorker,
		).
		WithController(
			webhookcachefiller.New(
				idpCache,
				installationNamespacePinnipedInformers.IDP().V1alpha1().WebhookIdentityProviders(),
				klogr.New(),
			),
			singletonWorker,
		).
		WithController(
			webhookcachecleaner.New(
				idpCache,
				installationNamespacePinnipedInformers.IDP().V1alpha1().WebhookIdentityProviders(),
				klogr.New(),
			),
			singletonWorker,
		)

	// Return a function which starts the informers and controllers.
	return func(ctx context.Context) {
		kubePublicNamespaceK8sInformers.Start(ctx.Done())
		installationNamespaceK8sInformers.Start(ctx.Done())
		installationNamespacePinnipedInformers.Start(ctx.Done())

		kubePublicNamespaceK8sInformers.WaitForCacheSync(ctx.Done())
		installationNamespaceK8sInformers.WaitForCacheSync(ctx.Done())
		installationNamespacePinnipedInformers.WaitForCacheSync(ctx.Done())

		go controllerManager.Start(ctx)
	}, nil
}

// Create the k8s clients that will be used by the controllers.
func createClients() (
	k8sClient *kubernetes.Clientset,
	aggregatorClient *aggregatorclient.Clientset,
	pinnipedClient *pinnipedclientset.Clientset,
	err error,
) {
	// Load the Kubernetes client configuration.
	kubeConfig, err := restclient.InClusterConfig()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not load in-cluster configuration: %w", err)
	}

	// explicitly use protobuf when talking to built-in kube APIs
	protoKubeConfig := createProtoKubeConfig(kubeConfig)

	// Connect to the core Kubernetes API.
	k8sClient, err = kubernetes.NewForConfig(protoKubeConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not initialize Kubernetes client: %w", err)
	}

	// Connect to the Kubernetes aggregation API.
	aggregatorClient, err = aggregatorclient.NewForConfig(protoKubeConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not initialize Kubernetes client: %w", err)
	}

	// Connect to the pinniped API.
	// I think we can't use protobuf encoding here because we are using CRDs
	// (for which protobuf encoding is not supported).
	pinnipedClient, err = pinnipedclientset.NewForConfig(kubeConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not initialize pinniped client: %w", err)
	}

	//nolint: nakedret
	return
}

// Create the informers that will be used by the controllers.
func createInformers(
	serverInstallationNamespace string,
	k8sClient *kubernetes.Clientset,
	pinnipedClient *pinnipedclientset.Clientset,
) (
	kubePublicNamespaceK8sInformers k8sinformers.SharedInformerFactory,
	installationNamespaceK8sInformers k8sinformers.SharedInformerFactory,
	installationNamespacePinnipedInformers pinnipedinformers.SharedInformerFactory,
) {
	kubePublicNamespaceK8sInformers = k8sinformers.NewSharedInformerFactoryWithOptions(
		k8sClient,
		defaultResyncInterval,
		k8sinformers.WithNamespace(issuerconfig.ClusterInfoNamespace),
	)
	installationNamespaceK8sInformers = k8sinformers.NewSharedInformerFactoryWithOptions(
		k8sClient,
		defaultResyncInterval,
		k8sinformers.WithNamespace(serverInstallationNamespace),
	)
	installationNamespacePinnipedInformers = pinnipedinformers.NewSharedInformerFactoryWithOptions(
		pinnipedClient,
		defaultResyncInterval,
		pinnipedinformers.WithNamespace(serverInstallationNamespace),
	)
	return
}

// Returns a copy of the input config with the ContentConfig set to use protobuf.
// Do not use this config to communicate with any CRD based APIs.
func createProtoKubeConfig(kubeConfig *restclient.Config) *restclient.Config {
	protoKubeConfig := restclient.CopyConfig(kubeConfig)
	const protoThenJSON = runtime.ContentTypeProtobuf + "," + runtime.ContentTypeJSON
	protoKubeConfig.AcceptContentTypes = protoThenJSON
	protoKubeConfig.ContentType = runtime.ContentTypeProtobuf
	return protoKubeConfig
}
