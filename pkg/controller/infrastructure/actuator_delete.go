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
	"encoding/hex"

	"github.com/23technologies/gardener-extension-provider-ionos/pkg/controller/infrastructure/ensurer"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis/transcoder"
	ionosapiwrapper "github.com/23technologies/ionos-api-wrapper/pkg"
	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
)

// delete deletes the infrastructure config.
//
// PARAMETERS
// ctx     context.Context                    Execution context
// infra   *extensionsv1alpha1.Infrastructure Infrastructure struct
// cluster *extensionscontroller.Cluster      Cluster struct
func (a *actuator) delete(ctx context.Context, infra *extensionsv1alpha1.Infrastructure, cluster *extensionscontroller.Cluster) error {
	actuatorConfig, err := a.getActuatorConfig(ctx, infra, cluster)
	// the shoot never reached the state of having a ProviderConfig assigned
	// we can assume nothing was setup
	if _, ok := err.(*transcoder.MissingProviderConfig); ok {
		return a.updateProviderStatus(ctx, infra, nil)
	}
	if err != nil {
		return err
	}

	client := ionosapiwrapper.GetClientForUser(actuatorConfig.user, actuatorConfig.password)

	infraStatus, _ := transcoder.DecodeInfrastructureStatusFromInfrastructure(infra)

	if nil != infraStatus {
		isDeleted := false

		labels, httpResponse, err := client.LabelApi.DatacentersLabelsGet(ctx, infraStatus.DatacenterID).Depth(1).Execute()
		if 404 == httpResponse.StatusCode {
			isDeleted = true
		} else {
			return err
		}

		if !isDeleted {
			for _, label := range *labels.Items {
				if nil != label.Properties && "cluster" == *label.Properties.Key {
					clusterValue, err := hex.DecodeString(*label.Properties.Value)
					if nil != err {
						continue
					}

					if infra.Namespace == string(clusterValue) {
						err = ensurer.EnsureDatacenterDeleted(ctx, client, infraStatus.DatacenterID)
						if err != nil {
							return err
						}

						isDeleted = true
						break
					}
				}
			}
		}

		if !isDeleted {
			err = ensurer.EnsureDHCPServerDeleted(ctx, client, infraStatus.DatacenterID, infraStatus.DHCPServerConfiguration.ServerID)
			if err != nil {
				return err
			}

			err = ensurer.EnsureNetworksDeleted(ctx, client, infraStatus.DatacenterID, infraStatus.NetworkIDs)
			if err != nil {
				return err
			}

			err = ensurer.EnsureFloatingPoolDeleted(ctx, client, infraStatus.FloatingPoolID)
			if err != nil {
				return err
			}
		}
	}

	return a.updateProviderStatus(ctx, infra, nil)
}
