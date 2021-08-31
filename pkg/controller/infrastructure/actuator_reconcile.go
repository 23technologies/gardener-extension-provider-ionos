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

	"github.com/23technologies/gardener-extension-provider-ionos/pkg/controller/infrastructure/ensurer"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis/transcoder"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis/v1alpha1"
	ionosapiwrapper "github.com/23technologies/ionos-api-wrapper/pkg"
	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// reconcile reconciles the infrastructure config.
//
// PARAMETERS
// ctx     context.Context                    Execution context
// infra   *extensionsv1alpha1.Infrastructure Infrastructure struct
// cluster *extensionscontroller.Cluster      Cluster struct
func (a *actuator) reconcile(ctx context.Context, infra *extensionsv1alpha1.Infrastructure, cluster *extensionscontroller.Cluster) error {
	actuatorConfig, err := a.getActuatorConfig(ctx, infra, cluster)
	if err != nil {
		return err
	}

	cpConfig, err := transcoder.DecodeControlPlaneConfigFromControllerCluster(cluster)
	if err != nil {
		return err
	}

	infraConfig, err := transcoder.DecodeInfrastructureConfigFromInfrastructure(infra)
	if err != nil {
		return err
	}

	infraStatus, _ := transcoder.DecodeInfrastructureStatusFromInfrastructure(infra)

	client := ionosapiwrapper.GetClientForUser(actuatorConfig.user, actuatorConfig.password)

	datacenterID, err := ensurer.EnsureDatacenter(ctx, client, cpConfig.Zone, infra.Namespace)
	if err != nil {
		return err
	}

	wanNetworkID, workerNetworkID, err := ensurer.EnsureNetworks(ctx, client, datacenterID, infra.Namespace, actuatorConfig.infraConfig.Networks)
	if err != nil {
		return err
	}

	dhcpServerID, err := ensurer.EnsureDHCPServer(ctx, client, datacenterID, cpConfig.Zone, actuatorConfig.infraConfig, infraStatus, workerNetworkID)
	if err != nil {
		return err
	}

	floatingPoolID, err := ensurer.EnsureFloatingPool(ctx, client, cpConfig.Zone, infra.Namespace, infraConfig.FloatingPool)
	if err != nil {
		return err
	}

	newInfraStatus := &v1alpha1.InfrastructureStatus{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
			Kind:       "InfrastructureStatus",
		},
		DatacenterID: datacenterID,
		DHCPServerConfiguration: &v1alpha1.DHCPServerConfigurationStatus{
			Cidr:     actuatorConfig.infraConfig.Networks.Workers,
			ServerID: dhcpServerID,
		},
		NetworkIDs: &v1alpha1.NetworkIDs{
			WAN: wanNetworkID,
		},
	}

	if "" != workerNetworkID {
		newInfraStatus.NetworkIDs.Workers = workerNetworkID
	}

	if "" != floatingPoolID {
		newInfraStatus.FloatingPoolID = floatingPoolID
	}

	return a.updateProviderStatus(ctx, infra, newInfraStatus)
}
