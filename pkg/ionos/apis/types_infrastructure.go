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

// InfrastructureConfig infrastructure configuration resource
type InfrastructureConfig struct {
	metav1.TypeMeta `json:",inline"`
	//
	DHCPServerConfiguration *DHCPServerConfiguration `json:"dhcpServerConfiguration"`

	// FloatingPool contains the floating pool configuration
	// +optional
	FloatingPool *FloatingPool `json:"floatingPool,omitempty"`
	// Networks contains the IONOS specific network configuration
	// +optional
	Networks *Networks `json:"networks,omitempty"`
}

// DHCPServerConfiguration contains the DHCP server configuration
type DHCPServerConfiguration struct {
	//
	Image    *MachineImage `json:"image"`
	//
	Cores    uint          `json:"cores"`
	//
	Memory   uint          `json:"memory"`
	//
	SSHKey   string        `json:"sshKey"`
	//
	UserData string        `json:"userData"`

	// Volume size
	VolumeSize float32 `json:"volumeSize,omitempty"`
}

// FloatingPool contains the floating pool configuration
type FloatingPool struct {
	// Floating pool name
	Name string `json:"name"`
	// Floating pool size
	Size uint16 `json:"size"`
}

// Networks contains the IONOS specific network configuration
type Networks struct {
	// Workers is a CIDRs of a worker subnet (private) to create (used for the VMs).
	Workers string `json:"workers"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InfrastructureStatus contains information about created infrastructure resources.
type InfrastructureStatus struct {
	metav1.TypeMeta `json:",inline"`
	// DatacenterID contains the IONOS data center ID.
	DatacenterID string `json:"datacenterID"`

	//
	// +optional
	DHCPServerConfiguration *DHCPServerConfigurationStatus `json:"dhcpServerConfiguration,omitempty"`
	// FloatingPoolID contains the floating pool IP ID.
	// +optional
	FloatingPoolID string `json:"floatingPoolID,omitempty"`
	// Networks is the ionos specific network configuration
	// +optional
	NetworkIDs *NetworkIDs `json:"networkIDs,omitempty"`
}

// DHCPServerConfiguration contains the DHCP server configuration
type DHCPServerConfigurationStatus struct {
	//
	Cidr     string `json:"cidr"`
	//
	ServerID string `json:"serverID"`
}

// Networks holds information about the Kubernetes and infrastructure networks.
type NetworkIDs struct {
	// WAN is the network ID for the public facing network interface.
	WAN string `json:"wan"`
	// Workers is the network ID of a worker subnet.
	Workers string `json:"workers"`
}
