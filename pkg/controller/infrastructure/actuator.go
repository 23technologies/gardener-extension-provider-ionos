/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package infrastructure contains functions used at the infrastructure controller
package infrastructure

import (
	"context"

	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis/controller"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis/transcoder"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis/v1alpha1"
	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/common"
	"github.com/gardener/gardener/extensions/pkg/controller/infrastructure"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type actuator struct {
	common.ChartRendererContext

	logger   logr.Logger
	gardenID string
}

type actuatorConfig struct {
	cloudProfileConfig *apis.CloudProfileConfig
	infraConfig        *apis.InfrastructureConfig
	user               string
	password           string
}

// NewActuator creates a new Actuator that updates the status of the handled Infrastructure resources.
func NewActuator(gardenID string) infrastructure.Actuator {
	return &actuator{
		logger:   log.Log.WithName("infrastructure-actuator"),
		gardenID: gardenID,
	}
}

func (a *actuator) getActuatorConfig(ctx context.Context, infra *extensionsv1alpha1.Infrastructure, cluster *extensionscontroller.Cluster) (*actuatorConfig, error) {
	cloudProfileConfig, err := transcoder.DecodeCloudProfileConfigFromControllerCluster(cluster)
	if err != nil {
		return nil, err
	}

	infraConfig, err := transcoder.DecodeInfrastructureConfigFromInfrastructure(infra)
	if err != nil {
		return nil, err
	}

	secret, err := extensionscontroller.GetSecretByReference(ctx, a.Client(), &infra.Spec.SecretRef)
	if err != nil {
		return nil, err
	}

	credentials, err := ionos.ExtractCredentials(secret)
	if err != nil {
		return nil, err
	}

	credentialsData := credentials.IonosCCM()

	config := &actuatorConfig{
		cloudProfileConfig: cloudProfileConfig,
		infraConfig:        infraConfig,
		user:               credentialsData.User,
		password:           credentialsData.Password,
	}

	return config, nil
}

// Delete implements infrastructure.Actuator.Delete
//
// PARAMETERS
// ctx     context.Context                    Execution context
// infra   *extensionsv1alpha1.Infrastructure Infrastructure struct
// cluster *extensionscontroller.Cluster      Cluster struct
func (a *actuator) Delete(ctx context.Context, infra *extensionsv1alpha1.Infrastructure, cluster *extensionscontroller.Cluster) error {
	return a.delete(ctx, infra, cluster)
}

// Migrate implements infrastructure.Actuator.Migrate
//
// PARAMETERS
// ctx     context.Context                    Execution context
// infra   *extensionsv1alpha1.Infrastructure Infrastructure struct
// cluster *extensionscontroller.Cluster      Cluster struct
func (a *actuator) Migrate(ctx context.Context, infra *extensionsv1alpha1.Infrastructure, cluster *extensionscontroller.Cluster) error {
	return nil
}

// Reconcile implements infrastructure.Actuator.Reconcile
//
// PARAMETERS
// ctx     context.Context                    Execution context
// infra   *extensionsv1alpha1.Infrastructure Infrastructure struct
// cluster *extensionscontroller.Cluster      Cluster struct
func (a *actuator) Reconcile(ctx context.Context, infra *extensionsv1alpha1.Infrastructure, cluster *extensionscontroller.Cluster) error {
	extendedCtx := context.WithValue(ctx, controller.CtxWrapDataKey("MethodData"), &controller.InfrastructureReconcileMethodData{})

	err := a.reconcile(extendedCtx, infra, cluster)

	if nil != err {
		a.reconcileOnErrorCleanup(extendedCtx, infra, cluster, err)
	}

	return err
}

// Restore implements infrastructure.Actuator.Restore
//
// PARAMETERS
// ctx     context.Context                    Execution context
// infra   *extensionsv1alpha1.Infrastructure Infrastructure struct
// cluster *extensionscontroller.Cluster      Cluster struct
func (a *actuator) Restore(ctx context.Context, infra *extensionsv1alpha1.Infrastructure, cluster *extensionscontroller.Cluster) error {
	return nil
}

// updateProviderStatus updates the infrastructure provider status.
//
// PARAMETERS
// ctx         context.Context                    Execution context
// infra       *extensionsv1alpha1.Infrastructure Infrastructure struct
// infraStatus *v1alpha1.InfrastructureStatus     Infrastructure status to be applied
func (a *actuator) updateProviderStatus(ctx context.Context, infra *extensionsv1alpha1.Infrastructure, infraStatus *v1alpha1.InfrastructureStatus) error {
	if nil == infraStatus {
		return nil
	}

	patch := client.MergeFrom(infra.DeepCopy())

	infra.Status.ProviderStatus = &runtime.RawExtension{
		Object: infraStatus,
	}

	return a.Client().Status().Patch(ctx, infra, patch)
}
