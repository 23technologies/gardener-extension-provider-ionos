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
	"net/http"
	"net/http/httptest"

	mockclient "github.com/gardener/gardener/pkg/mock/controller-runtime/client"
	mockkubernetes "github.com/gardener/gardener/pkg/client/kubernetes/mock"
	"github.com/golang/mock/gomock"
	ionossdk "github.com/ionos-cloud/sdk-go/v5"
	"github.com/onsi/ginkgo/v2"
)

const (
	TestDatacenterID = "01234567-89ab-4def-0123-c56789abcdef"
	TestFloatingPoolName = "MY-FLOATING-POOL"
	TestImageID = "01234567-89ab-4def-0123-c56789abcdef"
	TestImageName = "ubuntu"
	TestImageVersion = "1.0"
	TestNamespace = "test-namespace"
	TestRegion = "us"
	TestVolumeID = "3456789a-bcde-4012-3f56-789abcdef012"
	TestSSHPublicKey = "ecdsa-sha2-nistp384 AAAAE2VjZHNhLXNoYTItbmlzdHAzODQAAAAIbmlzdHAzODQAAABhBJ9S5cCzfygWEEVR+h3yDE83xKiTlc7S3pC3IadoYu/HAmjGPNRQZWLPCfZe5K3PjOGgXghmBY22voYl7bSVjy+8nZRPuVBuFDZJ9xKLPBImQcovQ1bMn8vXno4fvAF4KQ=="
	TestZone = "us-las"
)

// MockTestEnv represents the test environment for testing ionos API calls
type MockTestEnv struct {
	ChartApplier   *mockkubernetes.MockChartApplier
	Client         *mockclient.MockClient
	MockController *gomock.Controller
    StatusWriter   *mockclient.MockStatusWriter

	Server       *httptest.Server
	Mux          *http.ServeMux
	IonosClient  *ionossdk.APIClient
}

// Teardown shuts down the test environment
func (env *MockTestEnv) Teardown() {
	env.MockController.Finish()

	env.ChartApplier = nil
	env.Client = nil
	env.MockController = nil
	env.StatusWriter = nil

	env.Server.Close()

	env.Server = nil
	env.Mux = nil
	env.IonosClient = nil
}

// NewMockTestEnv generates a new, unconfigured test environment for testing purposes.
func NewMockTestEnv() MockTestEnv {
	ctrl := gomock.NewController(ginkgo.GinkgoT())

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	config := ionossdk.NewConfiguration("user", "dummy-password", "")

	config.Servers = ionossdk.ServerConfigurations{
		{
			URL: server.URL,
			Description: "Local mocked server base URL",
		},
	}

	ionosClient := ionossdk.NewAPIClient(config)

	return MockTestEnv{
		ChartApplier:   mockkubernetes.NewMockChartApplier(ctrl),
		Client:         mockclient.NewMockClient(ctrl),
		MockController: ctrl,
		StatusWriter:   mockclient.NewMockStatusWriter(ctrl),

		Server: server,
		Mux:    mux,
		IonosClient: ionosClient,
	}
}
