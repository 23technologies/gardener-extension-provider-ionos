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

// Package worker contains functions used at the worker controller
package worker

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos"

	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis/transcoder"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/worker"
	genericworkeractuator "github.com/gardener/gardener/extensions/pkg/controller/worker/genericactuator"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	mcmv1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachineClassKind yields the name of the machine class.
func (w *workerDelegate) MachineClassKind() string {
	return "MachineClass"
}

// MachineClass yields a newly initialized MachineClass object.
func (w *workerDelegate) MachineClass() client.Object {
	return &mcmv1alpha1.MachineClass{}
}

// MachineClassList yields a newly initialized MachineClassList object.
func (w *workerDelegate) MachineClassList() client.ObjectList {
	return &mcmv1alpha1.MachineClassList{}
}

// DeployMachineClasses generates and creates the ionos specific machine classes.
//
// PARAMETERS
// ctx context.Context Execution context
func (w *workerDelegate) DeployMachineClasses(ctx context.Context) error {
	if w.machineClasses == nil {
		if err := w.generateMachineConfig(ctx); err != nil {
			return err
		}
	}

	return w.seedChartApplier.Apply(ctx, filepath.Join(ionos.InternalChartsPath, "machineclass"), w.worker.Namespace, "machineclass", kubernetes.Values(map[string]interface{}{"machineClasses": w.machineClasses}))
}

// GenerateMachineDeployments generates the configuration for the desired machine deployments.
//
// PARAMETERS
// ctx context.Context Execution context
func (w *workerDelegate) GenerateMachineDeployments(ctx context.Context) (worker.MachineDeployments, error) {
	if w.machineDeployments == nil {
		if err := w.generateMachineConfig(ctx); err != nil {
			return nil, err
		}
	}
	return w.machineDeployments, nil
}

// getSecretData returns the secret referenced by the WorkerDelegate instance's spec.
//
// PARAMETERS
// ctx context.Context Execution context
func (w *workerDelegate) getSecretData(ctx context.Context) (*corev1.Secret, error) {
	return extensionscontroller.GetSecretByReference(ctx, w.Client(), &w.worker.Spec.SecretRef)
}

// generateMachineClassSecretData returns the machine class relevant secret values.
//
// PARAMETERS
// ctx context.Context Execution context
func (w *workerDelegate) generateMachineClassSecretData(ctx context.Context) (map[string][]byte, error) {
	secret, err := w.getSecretData(ctx)
	if err != nil {
		return nil, err
	}

	credentials, err := ionos.ExtractCredentials(secret)
	if err != nil {
		return nil, err
	}

	credentialsData := credentials.IonosMCM()

	return map[string][]byte{
		ionos.IonosCredentialsUserKey: []byte(credentialsData.User),
		ionos.IonosCredentialsPasswordKey: []byte(credentialsData.Password),
	}, nil
}

// generateMachineConfig generates the machine config of the WorkerDelegate instance's spec.
//
// PARAMETERS
// ctx context.Context Execution context
func (w *workerDelegate) generateMachineConfig(ctx context.Context) error {
	var (
		machineDeployments = worker.MachineDeployments{}
		machineClasses     []map[string]interface{}
		// machineImages      []apis.MachineImage
	)

	machineClassSecretData, err := w.generateMachineClassSecretData(ctx)
	if err != nil {
		return err
	}

	infraStatus, err := transcoder.DecodeInfrastructureStatusFromWorker(w.worker)
	if err != nil {
		return err
	}

	if len(w.worker.Spec.Pools) == 0 {
		return fmt.Errorf("missing pool")
	}

	for _, pool := range w.worker.Spec.Pools {
		workerPoolHash, err := worker.WorkerPoolHash(pool, w.cluster)
		if err != nil {
			return err
		}

		values, err := w.extractMachineValues(pool.MachineType)
		if err != nil {
			return errors.Wrap(err, "extracting machine values failed")
		}

		for _, zone := range pool.Zones {
			imageID, err := w.findMachineImageName(ctx, zone, pool.MachineImage.Name, pool.MachineImage.Version)
			if err != nil {
				return err
			}

			secretMap := map[string]interface{}{
				"userData": pool.UserData,
			}

			for key, value := range machineClassSecretData {
				secretMap[key] = value
			}

			machineClassSpec := map[string]interface{}{
				"datacenterID": infraStatus.DatacenterID,
				"cluster":      w.worker.Namespace,
				"zone":         zone,
				"cores":        values.Cores,
				"memory":       values.MemoryInMB,
				"imageID":      imageID,
				"sshKey":       string(w.worker.Spec.SSHPublicKey),
				"networkIDs":   infraStatus.NetworkIDs,
				"tags": map[string]string{
					"mcm.gardener.cloud/cluster": w.worker.Namespace,
					"mcm.gardener.cloud/role":    "node",
				},
				"secret": secretMap,
			}

			if "" != infraStatus.FloatingPoolID {
				machineClassSpec["floatingPoolID"] = infraStatus.FloatingPoolID
			}

			if values.MachineTypeOptions != nil {
				if len(values.MachineTypeOptions.ExtraConfig) > 0 {
					machineClassSpec["extraConfig"] = values.MachineTypeOptions.ExtraConfig
				}
			}

			if nil != pool.Volume && "" != pool.Volume.Size {
				volumeSizeQuantity, err := resource.ParseQuantity(pool.Volume.Size)
				if err != nil {
					return err
				}

				volumeSize, ok := volumeSizeQuantity.AsInt64()

				if ok && volumeSize > 0 {
					machineClassSpec["volumeSize"] = volumeSize
				}
			}

			deploymentName := fmt.Sprintf("%s-%s-%s", w.worker.Namespace, pool.Name, zone)
			className      := fmt.Sprintf("%s-%s", deploymentName, workerPoolHash)

			machineDeployments = append(machineDeployments, worker.MachineDeployment{
				Name:                 deploymentName,
				ClassName:            className,
				SecretName:           className,
				Minimum:              pool.Minimum,
				Maximum:              pool.Maximum,
				MaxSurge:             pool.MaxSurge,
				MaxUnavailable:       pool.MaxUnavailable,
				Labels:               pool.Labels,
				Annotations:          pool.Annotations,
				Taints:               pool.Taints,
				MachineConfiguration: genericworkeractuator.ReadMachineConfiguration(pool),
			})

			machineClassSpec["name"] = className

			machineClasses = append(machineClasses, machineClassSpec)
		}

	}
	w.machineDeployments = machineDeployments
	w.machineClasses = machineClasses

	return nil
}

type machineValues struct {
	Cores              int
	MemoryInMB         int
	MachineTypeOptions *apis.MachineTypeOptions
}

// extractMachineValues extracts the relevant machine values from the cloud profile spec.
//
// PARAMETERS
// ctx context.Context Execution context
func (w *workerDelegate) extractMachineValues(machineTypeName string) (*machineValues, error) {
	var machineType *corev1beta1.MachineType
	for _, mt := range w.cluster.CloudProfile.Spec.MachineTypes {
		if mt.Name == machineTypeName {
			machineType = &mt
			break
		}
	}
	if machineType == nil {
		err := fmt.Errorf("machine type %s not found in cloud profile spec", machineTypeName)
		return nil, err
	}

	values := &machineValues{}
	if n, ok := machineType.CPU.AsInt64(); ok {
		values.Cores = int(n)
	}
	if values.Cores <= 0 {
		err := fmt.Errorf("machine type %s has invalid CPU value %s", machineTypeName, machineType.CPU.String())
		return nil, err
	}

	if n, ok := machineType.Memory.AsInt64(); ok {
		values.MemoryInMB = int(n) / (1024 * 1024)
	}
	if values.MemoryInMB <= 0 {
		err := fmt.Errorf("machine type %s has invalid Memory value %s", machineTypeName, machineType.Memory.String())
		return nil, err
	}

	cloudProfileConfig, err := transcoder.DecodeConfigFromCloudProfile(w.cluster.CloudProfile)
	if err != nil {
		return nil, err
	}

	for _, mt := range cloudProfileConfig.MachineTypeOptions {
		if mt.Name == machineTypeName {
			values.MachineTypeOptions = &mt
			break
		}
	}

	return values, nil
}
