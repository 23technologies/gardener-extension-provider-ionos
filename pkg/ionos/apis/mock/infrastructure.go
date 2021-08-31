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

// Package mock provides all methods required to simulate a IONOS provider environment
package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis"
	"github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	TestInfrastructureName = "abc"
	TestInfrastructureProviderConfig = `{
		"apiVersion": "ionos.provider.extensions.gardener.cloud/v1alpha1",
		"kind": "InfrastructureConfig",
		"dhcpServerConfiguration": {
			"image": {
				"name": "ubuntu",
				"version": "20.04"
			},
			"sshKey": "ssh-rsa invalid",
			"userData": ""
		},
		"floatingPoolName": "MY-FLOATING-POOL",
		"networks": {"workers": "10.250.0.0/19"}
	}`
	TestInfrastructureSecretName = "cloudprovider"
	TestInfrastructureWorkersNetworkName = "test-namespace-workers"
)

// NewInfrastructure generates a new provider specification for testing purposes.
func NewInfrastructure() *v1alpha1.Infrastructure {
	return &v1alpha1.Infrastructure{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "extensions.gardener.cloud",
			Kind:       "Infrastructure",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      TestInfrastructureName,
			Namespace: TestNamespace,
		},
		Spec: v1alpha1.InfrastructureSpec{
			Region: TestRegion,
			SecretRef: corev1.SecretReference{
				Name: TestInfrastructureSecretName,
				Namespace: TestNamespace,
			},
			DefaultSpec: v1alpha1.DefaultSpec{
				ProviderConfig: &runtime.RawExtension{
					Raw: []byte(TestInfrastructureProviderConfig),
				},
			},
			SSHPublicKey: []byte(TestSSHPublicKey),
		},
	}
}

// NewInfrastructureConfigSpec generates a new infrastructure config specification for testing purposes.
func NewInfrastructureConfigSpec() *apis.InfrastructureConfig {
	return &apis.InfrastructureConfig{
		FloatingPool: &apis.FloatingPool{
			Name: TestFloatingPoolName,
			Size: 2,
		},
		Networks: &apis.Networks{
			Workers: TestInfrastructureWorkersNetworkName,
		},
	}
}

// ManipulateInfrastructure changes given provider specification.
//
// PARAMETERS
// infrastructure *extensions.Infrastructure Infrastructure specification
// data           map[string]interface{}     Members to change
func ManipulateInfrastructure(infrastructure *v1alpha1.Infrastructure, data map[string]interface{}) *v1alpha1.Infrastructure {
	for key, value := range data {
		if (strings.Index(key, "ObjectMeta") == 0) {
			manipulateStruct(&infrastructure.ObjectMeta, key[11:], value)
		} else if (strings.Index(key, "Spec") == 0) {
			manipulateStruct(&infrastructure.Spec, key[7:], value)
		} else if (strings.Index(key, "TypeMeta") == 0) {
			manipulateStruct(&infrastructure.TypeMeta, key[9:], value)
		} else {
			manipulateStruct(&infrastructure, key, value)
		}
	}

	return infrastructure
}

// SetupDatacentersEndpointOnMux configures a "/datacenters" endpoint on the mux given.
//
// PARAMETERS
// mux *http.ServeMux Mux to add handler to
func SetupDatacentersEndpointOnMux(mux *http.ServeMux) {
	mux.HandleFunc("/datacenters", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		if (strings.ToLower(req.Method) == "get") {
			res.WriteHeader(http.StatusOK)

			res.Write([]byte(fmt.Sprintf(`
{
	"id": %q,
	"type": "collection",
	"href": "",
	"items": [	]
}
			`, uuid.NewString())))
		} else if (strings.ToLower(req.Method) == "post") {
			res.WriteHeader(http.StatusCreated)

			jsonData := make([]byte, req.ContentLength)
			req.Body.Read(jsonData)

			var data map[string]interface{}

			jsonErr := json.Unmarshal(jsonData, &data)
			if jsonErr != nil {
				panic(jsonErr)
			}

			data["id"] = TestDatacenterID

			jsonData, jsonErr = json.Marshal(data)
			if jsonErr != nil {
				panic(jsonErr)
			}

			res.Write([]byte(jsonData))
		} else {
			panic("Unsupported HTTP method call")
		}
	})
}

// SetupTestDatacenterEndpointOnMux configures "/datacenters/<dcid>/*" endpoints on the mux given.
//
// PARAMETERS
// mux *http.ServeMux Mux to add handler to
func SetupTestDatacenterEndpointOnMux(mux *http.ServeMux) {
	baseURL := fmt.Sprintf("/datacenters/%s", TestDatacenterID)

	mux.HandleFunc(baseURL, func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		if (strings.ToLower(req.Method) == "get") {
			res.WriteHeader(http.StatusOK)

			res.Write([]byte(fmt.Sprintf(`
{
	"id": %q,
	"type": "datacenter",
	"href": "",
	"metadata": {
        "etag": "45480eb3fbfc31f1d916c1eaa4abdcc3",
        "createdDate": "2015-12-04T14:34:09.809Z",
        "createdBy": "user@example.com",
        "createdByUserId": "user@example.com",
        "lastModifiedDate": "2015-12-04T14:34:09.809Z",
        "lastModifiedBy": "user@example.com",
        "lastModifiedByUserId": "63cef532-26fe-4a64-a4e0-de7c8a506c90",
        "state": "AVAILABLE"
	},
	"properties": {
		"name": "Test 42",
		"description": "Proudly copied from the IONOS Cloud API documentation",
		"location": "us/las",
        "version": 8,
        "features": [],
        "secAuthProtection": true,
        "cpuArchitecture": {
            "cpuFamily": "AMD_OPTERON",
            "maxCores": "62",
            "maxRam": "245760",
            "vendor": "AuthenticAMD"
        }
	},
	"entities": {}
}
			`, TestDatacenterID)))
		} else {
			panic("Unsupported HTTP method call")
		}
	})

	mux.HandleFunc(fmt.Sprintf("%s/lans", baseURL), func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		if (strings.ToLower(req.Method) == "get") {
			res.WriteHeader(http.StatusOK)

			res.Write([]byte(fmt.Sprintf(`
{
	"id": %q,
	"type": "collection",
	"href": "",
	"items": [	]
}
			`, uuid.NewString())))
		} else if (strings.ToLower(req.Method) == "post") {
			res.WriteHeader(http.StatusCreated)

			jsonData := make([]byte, req.ContentLength)
			req.Body.Read(jsonData)

			var data map[string]interface{}

			jsonErr := json.Unmarshal(jsonData, &data)
			if jsonErr != nil {
				panic(jsonErr)
			}

			data["id"] = uuid.NewString()

			jsonData, jsonErr = json.Marshal(data)
			if jsonErr != nil {
				panic(jsonErr)
			}

			res.Write([]byte(jsonData))
		} else {
			panic("Unsupported HTTP method call")
		}
	})
}
