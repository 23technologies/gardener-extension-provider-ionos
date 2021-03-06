/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 5.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// KubernetesConfigProperties struct for KubernetesConfigProperties
type KubernetesConfigProperties struct {
	// A Kubernetes Config file data
	Kubeconfig *string `json:"kubeconfig,omitempty"`
}



// GetKubeconfig returns the Kubeconfig field value
// If the value is explicit nil, the zero value for string will be returned
func (o *KubernetesConfigProperties) GetKubeconfig() *string {
	if o == nil {
		return nil
	}


	return o.Kubeconfig

}

// GetKubeconfigOk returns a tuple with the Kubeconfig field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesConfigProperties) GetKubeconfigOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Kubeconfig, true
}

// SetKubeconfig sets field value
func (o *KubernetesConfigProperties) SetKubeconfig(v string) {


	o.Kubeconfig = &v

}

// HasKubeconfig returns a boolean if a field has been set.
func (o *KubernetesConfigProperties) HasKubeconfig() bool {
	if o != nil && o.Kubeconfig != nil {
		return true
	}

	return false
}

func (o KubernetesConfigProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Kubeconfig != nil {
		toSerialize["kubeconfig"] = o.Kubeconfig
	}
	return json.Marshal(toSerialize)
}

type NullableKubernetesConfigProperties struct {
	value *KubernetesConfigProperties
	isSet bool
}

func (v NullableKubernetesConfigProperties) Get() *KubernetesConfigProperties {
	return v.value
}

func (v *NullableKubernetesConfigProperties) Set(val *KubernetesConfigProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableKubernetesConfigProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableKubernetesConfigProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKubernetesConfigProperties(val *KubernetesConfigProperties) *NullableKubernetesConfigProperties {
	return &NullableKubernetesConfigProperties{value: val, isSet: true}
}

func (v NullableKubernetesConfigProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKubernetesConfigProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


