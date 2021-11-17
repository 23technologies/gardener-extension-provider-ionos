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

// Package ensurer provides functions used to ensure infrastructure changes to be applied
package ensurer

import (
	"context"
	"fmt"
	"strings"

	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis"
	ionossdk "github.com/ionos-cloud/sdk-go/v5"
)

func createLANNetwork(ctx context.Context, client *ionossdk.APIClient, datacenterID, name string, public bool) (string, error) {
	lanProperties := ionossdk.LanPropertiesPost{
		Name: &name,
		Public: &public,
	}

	lanApiCreateRequest := client.LanApi.DatacentersLansPost(ctx, datacenterID).Depth(0)
	lan, _, err := lanApiCreateRequest.Lan(ionossdk.LanPost{Properties: &lanProperties}).Execute()
	if nil != err {
		return "", err
	}

	return *lan.Id, nil
}

func EnsureFloatingPool(ctx context.Context, client *ionossdk.APIClient, zone, namespace string, floatingPool *apis.FloatingPool) (string, error) {
	if nil != floatingPool {
		name := fmt.Sprintf("%s-%s", namespace, floatingPool.Name)

		ipBlocks, _, err := client.IPBlocksApi.IpblocksGet(ctx).Depth(1).Execute()
		if nil != err {
			return "", err
		}

		for _, ipBlock := range *ipBlocks.Items {
			if name == *ipBlock.Properties.Name {
				return *ipBlock.Id, nil
			}
		}

		location := strings.Replace(zone, "-", "/", 1)
		ipBlockSize := int32(floatingPool.Size)

		ipBlockProperties := ionossdk.IpBlockProperties{
			Name:     &name,
			Location: &location,
			Size:     &ipBlockSize,
		}

		ipBlockApiCreateRequest := client.IPBlocksApi.IpblocksPost(ctx).Depth(0)
		ipBlock, _, err := ipBlockApiCreateRequest.Ipblock(ionossdk.IpBlock{Properties: &ipBlockProperties}).Execute()
		if nil != err {
			return "", err
		}

		return *ipBlock.Id, nil
	}

	return "", nil
}

func EnsureFloatingPoolDeleted(ctx context.Context, client *ionossdk.APIClient, floatingPoolID string) error {
	if "" != floatingPoolID {
		_, httpResponse, err := client.IPBlocksApi.IpblocksDelete(ctx, floatingPoolID).Depth(0).Execute()
		if (nil != err) && (404 != httpResponse.StatusCode) {
			return err
		}
	}

	return nil
}

// EnsureNetworks verifies the network resources requested are available.
//
// PARAMETERS
// ctx       context.Context Execution context
// client    *hcloud.Client  HCloud client
// namespace string          Shoot namespace
// networks  *apis.Networks  Networks struct
func EnsureNetworks(ctx context.Context, client *ionossdk.APIClient, datacenterID, namespace string, networks *apis.Networks) (string, string, error) {
	var (
		wanID string
		workersID string
	)

	name := fmt.Sprintf("%s-wan", namespace)

	lans, _, err := client.LanApi.DatacentersLansGet(ctx, datacenterID).Depth(1).Execute()
	if nil != err {
		return wanID, workersID, err
	}

	for _, lan := range *lans.Items {
		if name == *lan.Properties.Name {
			wanID = *lan.Id
		}
	}

	if "" == wanID {
		wanID, err = createLANNetwork(ctx, client, datacenterID, name, true)
		if nil != err {
			return wanID, workersID, err
		}
	}

	if "" != networks.Workers {
		name = fmt.Sprintf("%s-workers", namespace)

		for _, lan := range *lans.Items {
			if name == *lan.Properties.Name {
				workersID = *lan.Id
			}
		}

		if "" == workersID {
			workersID, err = createLANNetwork(ctx, client, datacenterID, name, false)
			if nil != err {
				return wanID, workersID, err
			}
		}
	}

	return wanID, workersID, nil
}

// EnsureNetworksDeleted removes any previously created network resources.
//
// PARAMETERS
// ctx       context.Context  Execution context
// client    *hcloud.Client   HCloud client
// namespace string           Shoot namespace
// networks  *apis.NetworkIDs Network IDs struct
func EnsureNetworksDeleted(ctx context.Context, client *ionossdk.APIClient, datacenterID string, networkIDs *apis.NetworkIDs) error {
	if "" != networkIDs.WAN {
		_, httpResponse, err := client.LanApi.DatacentersLansDelete(ctx, datacenterID, networkIDs.WAN).Depth(0).Execute()
		if (nil != err) && (404 != httpResponse.StatusCode) {
			return err
		}
	}

	if "" != networkIDs.Workers {
		_, httpResponse, err := client.LanApi.DatacentersLansDelete(ctx, datacenterID, networkIDs.Workers).Depth(0).Execute()
		if (nil != err) && (404 != httpResponse.StatusCode) {
			return err
		}
	}

	return nil
}
