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

// Package apis is the main package for ionos specific APIs
package apis

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudProfileConfig contains provider-specific configuration that is embedded into Gardener's `CloudProfile`
// resource.
type CloudProfileConfig struct {
	metav1.TypeMeta
	// Regions is the specification of regions and zones topology
	Regions []RegionSpec
	// MachineImages is the list of machine images that are understood by the controller. It maps
	// logical names and versions to provider-specific identifiers.
	MachineImages []MachineImages
	// DefaultClassStoragePolicyName is the name of the ionos storage policy to use for the 'default-class' storage class
	DefaultClassStoragePolicyName string
	// MachineTypeOptions is the list of machine type options to set additional options for individual machine types.
	MachineTypeOptions []MachineTypeOptions
	// DockerDaemonOptions contains configuration options for docker daemon service
	DockerDaemonOptions *DockerDaemonOptions
}

// RegionSpec specifies the topology of a region and its zones.
// A region consists of a Vcenter host, transport zone and optionally a data center.
// A zone in a region consists of a data center (if not specified in the region), a computer cluster,
// and optionally a resource zone or host system.
type RegionSpec struct {
	// Name is the name of the region
	Name string

	// MachineImages is the list of machine images that are understood by the controller. If provided, it overwrites the global
	// MachineImages of the CloudProfileConfig
	MachineImages []MachineImages
}

// MachineImages is a mapping from logical names and versions to provider-specific identifiers.
type MachineImages struct {
	// Name is the logical name of the machine image.
	Name string
	// Versions contains versions and a provider-specific identifier.
	Versions []MachineImageVersion
}

// MachineImageVersion contains a version and a provider-specific identifier.
type MachineImageVersion struct {
	// Version is the version of the image.
	Version string

	// ImageName is the IONOS Cloud image name if not matching name + "-" + version.
	// +optional
	ImageName string `json:"imageName,omitempty"`
}

// MachineTypeOptions defines additional VM options for an machine type given by name
type MachineTypeOptions struct {
	// Name is the name of the machine type
	Name string
	// ExtraConfig allows to specify additional VM options.
	// e.g. sched.swap.vmxSwapEnabled=false to disable the VMX process swap file
	ExtraConfig map[string]string
}

// DockerDaemonOptions contains configuration options for Docker daemon service
type DockerDaemonOptions struct {
	// HTTPProxyConf contains HTTP/HTTPS proxy configuration for Docker daemon
	HTTPProxyConf *string
	// InsecureRegistries adds the given registries to Docker on the worker nodes
	// (see https://docs.docker.com/registry/insecure/)
	InsecureRegistries []string
}
