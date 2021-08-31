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

// Package ionos provides types and functions used for ionos interaction
package ionos

import (
	"path/filepath"

	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
)

const (
	// Name is the name of the IONOS provider controller.
	Name = "provider-ionos"

	// MachineControllerManagerImageName is the name of the MachineControllerManager image.
	MachineControllerManagerImageName = "machine-controller-manager"
	// MCMProviderIonosImageName is the namne of the IONOS provider plugin image.
	MCMProviderIonosImageName = "machine-controller-manager-provider-ionos"
	// CloudControllerImageName is the name of the external IONOS CloudProvider image.
	CloudControllerImageName = "ionos-cloud-controller-manager"

	// CloudControllerManagerName is the constant for the name of the CloudController deployed by the control plane controller.
	CloudControllerManagerName = "cloud-controller-manager"

	// CloudControllerManagerServerName is the constant for the name of the CloudController deployed by the control plane controller.
	CloudControllerManagerServerName = "cloud-controller-manager-server"
	// CSIAttacherImageName is the name of the CSI attacher image.
	CSIAttacherImageName = "csi-attacher"
	// CSINodeDriverRegistrarImageName is the name of the CSI driver registrar image.
	CSINodeDriverRegistrarImageName = "csi-node-driver-registrar"
	// CSIProvisionerImageName is the name of the CSI provisioner image.
	CSIProvisionerImageName = "csi-provisioner"
	// CSIDriverControllerImageName is the name of the CSI driver controller plugin image.
	CSIDriverControllerImageName = "csi-driver-controller"
	// CSIDriverNodeImageName is the name of the CSI driver node plugin image.
	CSIDriverNodeImageName = "csi-driver-node"
	// CSIResizerImageName is the name of the csi-resizer image.
	CSIResizerImageName = "csi-resizer"
	// LivenessProbeImageName is the name of the liveness-probe image.
	LivenessProbeImageName = "liveness-probe"

	// Cloud Controller Manager IONOS credentials user name key
	IonosCCMCredentialsUserKey     = "ionosCCMUser"
	// Cloud Controller Manager IONOS credentials password key
	IonosCCMCredentialsPasswordKey = "ionosCCMPassword"
	// Common IONOS credentials user name key
	IonosCredentialsUserKey     = "ionosUser"
	// Common IONOS credentials password key
	IonosCredentialsPasswordKey = "ionosPassword"
	// Container Storage Interface driver IONOS credentials user name key
	IonosCSICredentialsUserKey     = "ionosCSIUser"
	// Container Storage Interface driver IONOS credentials password key
	IonosCSICredentialsPasswordKey = "ionosCSIPassword"
	// Machine Controller Manager IONOS credentials user name key
	IonosMCMCredentialsUserKey     = "ionosMCMUser"
	// Machine Controller Manager IONOS credentials password key
	IonosMCMCredentialsPasswordKey = "ionosMCMPassword"

	// CloudProviderConfig is the name of the configmap containing the cloud provider config.
	CloudProviderConfig = "cloud-provider-config"
	// CloudProviderConfigMapKey is the key storing the cloud provider config as value in the cloud provider configmap.
	CloudProviderConfigMapKey = "cloudprovider.conf"
	// MachineControllerManagerName is a constant for the name of the machine-controller-manager.
	MachineControllerManagerName = "machine-controller-manager"
	// MachineControllerManagerVpaName is the name of the VerticalPodAutoscaler of the machine-controller-manager deployment.
	MachineControllerManagerVpaName = "machine-controller-manager-vpa"
	// MachineControllerManagerMonitoringConfigName is the name of the ConfigMap containing monitoring stack configurations for machine-controller-manager.
	MachineControllerManagerMonitoringConfigName = "machine-controller-manager-monitoring-config"

	// CSIProvisionerName is a constant for the name of the csi-provisioner component.
	CSIProvisionerName = "csi-provisioner"
	// CSIAttacherName is a constant for the name of the csi-attacher component.
	CSIAttacherName = "csi-attacher"
	// CSIResizerName is a constant for the name of the csi-resizer component.
	CSIResizerName = "csi-resizer"
	// CSIControllerName is a constant for the name of the ionos-csi-controller component.
	CSIControllerName = "ionos-csi-controller"
	// CSINodeName is a constant for the chart name for a CSI node deployment in the shoot.
	CSINodeName = "ionos-csi-node"
	// CSIDriverName is a constant for the name of the csi-driver component.
	CSIDriverName = "csi-driver"
)

var (
	// ChartsPath is the path to the charts
	ChartsPath = filepath.Join("charts")
	// InternalChartsPath is the path to the internal charts
	InternalChartsPath = filepath.Join(ChartsPath, "internal")

	// UsernamePrefix is a constant for the username prefix of components deployed by OpenStack.
	UsernamePrefix = extensionsv1alpha1.SchemeGroupVersion.Group + ":" + Name + ":"
)
