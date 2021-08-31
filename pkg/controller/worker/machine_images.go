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

	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis/transcoder"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis/v1alpha1"
	ionosapiwrapper "github.com/23technologies/ionos-api-wrapper/pkg"
	"github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/worker"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"
)

// findMachineImageName returns the image name for the given name and version values.
//
// PARAMETERS
// ctx     context.Context Execution context
// name    string          Machine image name
// version string          Machine image version
func (w *workerDelegate) findMachineImageName(ctx context.Context, zone, name, version string) (string, error) {
	secret, err := w.getSecretData(ctx)
	if nil != err {
		return "", err
	}

	credentials, err := ionos.ExtractCredentials(secret)
	if nil != err {
		return "", err
	}

	customImageNameVariant := fmt.Sprintf("%s-%s", name, version)
	client := ionosapiwrapper.GetClientForUser(credentials.IonosMCM().User, credentials.IonosMCM().Password)

	imageID, err := apis.FindMachineImageName(ctx, client, zone, name, version, customImageNameVariant)
	if nil != err {
		return "", worker.ErrorMachineImageNotFound(name, version)
	}

	return imageID, nil
}

// UpdateMachineImagesStatus adds machineImages to the `WorkerStatus` resource.
//
// PARAMETERS
// ctx context.Context Execution context
func (w *workerDelegate) UpdateMachineImagesStatus(ctx context.Context) error {
	if w.machineImages == nil {
		if err := w.generateMachineConfig(ctx); nil != err {
			return err
		}
	}

	var workerStatus *apis.WorkerStatus
	var workerStatusV1alpha1 *v1alpha1.WorkerStatus

	if w.worker.Status.ProviderStatus == nil {
		workerStatus = &apis.WorkerStatus{
			TypeMeta: metav1.TypeMeta{
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
				Kind:       "WorkerStatus",
			},
			MachineImages: w.machineImages,
		}
	} else {
		// Decode the current worker provider status.
		decodedWorkerStatus, err := transcoder.DecodeWorkerStatusFromWorker(w.worker)
		if nil != err {
			return err
		}

		workerStatus = decodedWorkerStatus
		workerStatus.MachineImages = w.machineImages

		workerStatusV1alpha1 = &v1alpha1.WorkerStatus{
			TypeMeta: metav1.TypeMeta{
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
				Kind:       "WorkerStatus",
			},
		}
	}

	workerStatusV1alpha1 = &v1alpha1.WorkerStatus{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
			Kind:       "WorkerStatus",
		},
	}

	if err := w.Scheme().Convert(workerStatus, workerStatusV1alpha1, nil); nil != err {
		return err
	}

	return controller.TryUpdateStatus(ctx, retry.DefaultBackoff, w.Client(), w.worker, func() error {
		w.worker.Status.ProviderStatus = &runtime.RawExtension{Object: workerStatusV1alpha1}
		return nil
	})
}
