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
	"encoding/hex"
	"fmt"
	"strings"

	ionosapiwrapper "github.com/23technologies/ionos-api-wrapper/pkg"
	ionossdk "github.com/ionos-cloud/sdk-go/v5"
)

func EnsureDatacenter(ctx context.Context, client *ionossdk.APIClient, zone, namespace string) (string, error) {
	name := fmt.Sprintf("%s-%s", namespace, zone)

	datacenters, _, err := client.DataCenterApi.DatacentersGet(ctx).Depth(1).Execute()
	if nil != err {
		return "", err
	}

	for _, datacenter := range *datacenters.Items {
		if name == *datacenter.Properties.Name {
			return *datacenter.Id, nil
		}
	}

	location := strings.Replace(zone, "-", "/", 1)

	datacenterProperties := ionossdk.DatacenterProperties{
		Name: &name,
		Location: &location,
	}

	datacenterApiCreateRequest := client.DataCenterApi.DatacentersPost(ctx).Depth(0)
	datacenter, _, err := datacenterApiCreateRequest.Datacenter(ionossdk.Datacenter{Properties: &datacenterProperties}).Execute()
	if nil != err {
		return "", err
	}

	datacenterID := *datacenter.Id

	err = ionosapiwrapper.WaitForDatacenterModifications(ctx, client, datacenterID)
	if nil != err {
		return "", err
	}

	err = ionosapiwrapper.AddLabelToDatacenter(ctx, client, datacenterID, "cluster", hex.EncodeToString([]byte(namespace)))
	if nil != err {
		return "", err
	}

	return datacenterID, nil
}

func EnsureDatacenterDeleted(ctx context.Context, client *ionossdk.APIClient, datacenterID string) error {
	_, httpResponse, err := client.DataCenterApi.DatacentersDelete(ctx, datacenterID).Depth(0).Execute()
	if (nil != err) && (404 != httpResponse.StatusCode) {
		return err
	}

	return nil
}
