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
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis"
	"github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	jsonInfrastructureImageTemplate = `{
		"id": %q,
		"type": "image",
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
			"name": "%s-%s",
			"description": "Proudly copied from the IONOS Cloud API documentation",
			"location": "us/las",
			"size": 100,
			"cpuHotPlug": true,
			"cpuHotUnplug": true,
			"ramHotPlug": true,
			"ramHotUnplug": true,
			"nicHotPlug": true,
			"nicHotUnplug": true,
			"discVirtioHotPlug": true,
			"discVirtioHotUnplug": true,
			"discScsiHotPlug": true,
			"discScsiHotUnplug": true,
			"licenceType": "LINUX",
			"imageType": "HDD",
			"public": true,
			"imageAliases": [],
			"cloudInit": "V1"
		}
	}`
	jsonInfrastructureNicTemplate = `{
		"id": %q,
		"type": "nic",
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
			"name": "NIC",
			"mac": "00:11:22:33:44:55",
			"ips": [],
			"dhcp": true,
			"lan": 2,
			"firewallActive": false,
			"firewallType": "INGRESS",
			"deviceNumber": 1,
			"pciSlot": 1
		},
		"entities": {
			"flowlogs": {},
			"firewallrules": {}
		}
	}
	`
	jsonInfrastructureServerDataTemplate = `
{
	"id": %q,
	"type": "server",
	"href": "",
	"metadata": {
		"etag": "45480eb3fbfc31f1d916c1eaa4abdcc3",
		"createdDate": "2015-12-04T14:34:09.809Z",
		"createdBy": "user@example.com",
		"createdByUserId": "user@example.com",
		"lastModifiedDate": "2015-12-04T14:34:09.809Z",
		"lastModifiedBy": "user@example.com",
		"lastModifiedByUserId": "63cef532-26fe-4a64-a4e0-de7c8a506c90",
		"state": %q
	},
	"properties": {
		"templateUuid": "15f67991-0f51-4efc-a8ad-ef1fb31a480c",
		"name": %q,
		"cores": 4,
		"ram": 4096,
		"availabilityZone": "AUTO",
		"vmState": %q,
		"bootCdrom": {
			"id": "",
			"type": "resource",
			"href": ""
		},
		"bootVolume": {
			"id": %q,
			"type": "resource",
			"href": ""
		},
		"cpuFamily": "AMD_OPTERON",
		"type": "CUBE"
	},
	"entities": {
		"cdroms": {},
		"volumes": {
			"id": "15f67991-0f51-4efc-a8ad-ef1fb31a480c",
			"type": "collection",
			"href": "",
			"items": [%s],
			"offset": 0,
			"limit": 1000,
			"_links": {}
		},
		"nics": {}
	}
}
	`
	jsonVolumeTemplate = `
{
	"id": %q,
	"type": "volume",
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
		"name": "My resource",
		"type": "HDD",
		"size": 100,
		"availabilityZone": "AUTO",
		"image": "15f67991-0f51-4efc-a8ad-ef1fb31a480c",
		"imagePassword": null,
		"sshKeys": [],
		"bus": "VIRTIO",
		"licenceType": "LINUX",
		"cpuHotPlug": true,
		"ramHotPlug": true,
		"nicHotPlug": true,
		"nicHotUnplug": true,
		"discVirtioHotPlug": true,
		"discVirtioHotUnplug": true,
		"deviceNumber": 3,
		"pciSlot": 7,
		"userData": ""
	}
}
	`
	TestInfrastructureName = "abc"
	TestInfrastructureProviderConfig = `{
		"apiVersion": "ionos.provider.extensions.gardener.cloud/v1alpha1",
		"kind": "InfrastructureConfig",
		"dhcpServerConfiguration": {
			"image": {
				"name": "test",
				"version": "1.0"
			},
			"ip": "10.250.31.254",
			"userData": "#cloud-config"
		},
		"floatingPoolName": "MY-FLOATING-POOL",
		"networks": {"workers": "10.250.0.0/19"}
	}`
	TestInfrastructureSecretName = "cloudprovider"
	TestInfrastructureServerID = "6789abcd-ef01-4345-6789-abcdef012325"
	TestInfrastructureServerNameTemplate = "machine-%s"
	TestInfrastructureServerNicID = "23456789-abcd-4f01-23e5-6789abcdef01"
	TestInfrastructureWorkersNetworkName = "test-namespace-workers"
)

var testInfrastructureNicID *int32 = new(int32)

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

// newJsonServerData generates a JSON server data object for testing purposes.
//
// PARAMETERS
// serverID    string Server ID to use
// serverState string Server state to use
func newJsonServerData(serverID string, serverState string) string {
	serverBootState := "RUNNING"

	if "AVAILABLE" != serverState {
		serverBootState = "SHUTOFF"
	}

	jsonVolumeData := fmt.Sprintf(jsonVolumeTemplate, TestVolumeID)
	TestInfrastructureServerName := fmt.Sprintf(TestInfrastructureServerNameTemplate, serverID)
	return fmt.Sprintf(jsonInfrastructureServerDataTemplate, serverID, serverState, TestInfrastructureServerName, serverBootState, TestVolumeID, jsonVolumeData)
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

// SetupImagesEndpointOnMux configures a "/images" endpoint on the mux given.
//
// PARAMETERS
// mux *http.ServeMux Mux to add handler to
func SetupImagesEndpointOnMux(mux *http.ServeMux) {
	mux.HandleFunc("/images", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		res.WriteHeader(http.StatusOK)

		res.Write([]byte(fmt.Sprintf(`
{
	"id": %q,
	"type": "collection",
	"href": "",
	"items": [
		%s
	]
}
		`, uuid.NewString(), fmt.Sprintf(jsonInfrastructureImageTemplate, TestImageID, TestImageName, TestImageVersion))))
	})
}

// SetupServersEndpointOnMux configures a "/datacenters/<id>/servers" endpoint on the mux given.
//
// PARAMETERS
// mux *http.ServeMux Mux to add handler to
func SetupServersEndpointOnMux(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("/datacenters/%s/servers", TestDatacenterID), func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		if (strings.ToLower(req.Method) == "get") {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(fmt.Sprintf(`
{
	"id": %q,
	"type": "collection",
	"href": "",
	"items": [
		%s
	]
}
			`, uuid.NewString(), newJsonServerData(TestInfrastructureServerID, "AVAILABLE"))))

		} else if (strings.ToLower(req.Method) == "post") {
			res.WriteHeader(http.StatusAccepted)

			jsonData := make([]byte, req.ContentLength)
			req.Body.Read(jsonData)

			var data map[string]interface{}

			jsonErr := json.Unmarshal(jsonData, &data)
			if jsonErr != nil {
				panic(jsonErr)
			}

			res.Write([]byte(newJsonServerData(TestInfrastructureServerID, "BUSY")))
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

	mux.HandleFunc(fmt.Sprintf("%s/labels", baseURL), func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		if (strings.ToLower(req.Method) == "get") {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(fmt.Sprintf(`
{
"id": %q,
"type": "collection",
"href": "",
"items": [ ]
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

			res.Write([]byte(jsonData))
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

			data["id"] = strconv.Itoa(int(*testInfrastructureNicID))

			jsonData, jsonErr = json.Marshal(data)
			if jsonErr != nil {
				panic(jsonErr)
			}

			res.Write([]byte(jsonData))
			atomic.AddInt32(testInfrastructureNicID, 1)
		} else {
			panic("Unsupported HTTP method call")
		}
	})

	mux.HandleFunc(fmt.Sprintf("%s/volumes", baseURL), func(res http.ResponseWriter, req *http.Request) {
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

			res.Write([]byte(fmt.Sprintf(jsonVolumeTemplate, TestVolumeID)))
		} else {
			panic("Unsupported HTTP method call")
		}
	})

	mux.HandleFunc(fmt.Sprintf("%s/volumes/%s", baseURL, TestVolumeID), func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		if (strings.ToLower(req.Method) == "get") {
			res.WriteHeader(http.StatusOK)

			res.Write([]byte(fmt.Sprintf(jsonVolumeTemplate, TestVolumeID)))
		} else {
			panic("Unsupported HTTP method call")
		}
	})
}

// SetupTestInfrastructureServerEndpointOnMux configures "/datacenters/<dcid>/servers/<sid>/*" endpoints on the mux given.
//
// PARAMETERS
// mux *http.ServeMux Mux to add handler to
func SetupTestInfrastructureServerEndpointOnMux(mux *http.ServeMux) {
	baseURL := fmt.Sprintf("/datacenters/%s/servers/%s", TestDatacenterID, TestInfrastructureServerID)

	mux.HandleFunc(baseURL, func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		if (strings.ToLower(req.Method) == "delete") {
			res.WriteHeader(http.StatusAccepted)
		} else if (strings.ToLower(req.Method) == "get") {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(newJsonServerData(TestInfrastructureServerID, "AVAILABLE")))
		} else {
			panic("Unsupported HTTP method call")
		}
	})

	mux.HandleFunc(fmt.Sprintf("%s/nics", baseURL), func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		if (strings.ToLower(req.Method) == "post") {
			res.WriteHeader(http.StatusAccepted)
			res.Write([]byte(fmt.Sprintf(jsonInfrastructureNicTemplate, TestInfrastructureServerNicID)))
		} else {
			panic("Unsupported HTTP method call")
		}
	})

	mux.HandleFunc(fmt.Sprintf("%s/nics/%s", baseURL, TestInfrastructureServerNicID), func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		if (strings.ToLower(req.Method) == "get") {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(fmt.Sprintf(jsonInfrastructureNicTemplate, TestInfrastructureServerNicID)))
		} else {
			panic("Unsupported HTTP method call")
		}
	})

	mux.HandleFunc(fmt.Sprintf("%s/start", baseURL), func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		if (strings.ToLower(req.Method) == "post") {
			res.WriteHeader(http.StatusAccepted)
		} else {
			panic("Unsupported HTTP method call")
		}
	})

	mux.HandleFunc(fmt.Sprintf("%s/stop", baseURL), func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		if (strings.ToLower(req.Method) == "post") {
			res.WriteHeader(http.StatusAccepted)
		} else {
			panic("Unsupported HTTP method call")
		}
	})
}

// SetupTestImageEndpointOnMux configures a "/images" endpoint on the mux given.
//
// PARAMETERS
// mux *http.ServeMux Mux to add handler to
func SetupTestImageEndpointOnMux(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("/images/%s", TestImageID), func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json; charset=utf-8")

		res.WriteHeader(http.StatusOK)

		res.Write([]byte(fmt.Sprintf(jsonInfrastructureImageTemplate, TestImageID, TestImageName, TestImageVersion)))
	})
}
