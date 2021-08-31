// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package apis is the main package for ionos specific APIs
package apis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	ionossdk "github.com/ionos-cloud/sdk-go/v5"
)

// FindMachineImageName returns the image name for the given name and version values.
//
// PARAMETERS
// ctx     context.Context Execution context
// name    string          Machine image name
// version string          Machine image version
func FindMachineImageName(ctx context.Context, client *ionossdk.APIClient, zone, name, version, customImageNameVariant string) (string, error) {
	images, _, err := client.ImageApi.ImagesGet(ctx).Depth(1).Execute()
	if nil != err {
		return "", err
	}

	if "" == customImageNameVariant {
		customImageNameVariant = fmt.Sprintf("%s-%s", name, version)
	}

	defaultImageNameVariant := fmt.Sprintf("%s-%s.qcow2", name, version)
	defaultImageNameVariantWithCloudInit := fmt.Sprintf("%s-%s-cloud-init.qcow2", name, version)

	location := strings.Replace(zone, "-", "/", 1)

	for _, image := range *images.Items {
		imageName := strings.ToLower(*image.Properties.Name)

		if "AVAILABLE" != *image.Metadata.State {
			continue
		} else if location != *image.Properties.Location {
			continue
		} else if (!image.Properties.HasCloudInit() || "NONE" == *image.Properties.CloudInit) {
			continue
		} else if customImageNameVariant != imageName && defaultImageNameVariant != imageName && defaultImageNameVariantWithCloudInit != imageName {
			continue
		}

		return *image.Id, nil
	}

	return "", errors.New(fmt.Sprintf("No matching image found for %s with version %s", name, version))
}

// GetRegionFromZone returns the region for a given zone string
//
// PARAMETERS
// zone string Zone
func GetRegionFromZone(zone string) string {
	zoneData := strings.SplitN(zone, "-", 2)
	return zoneData[0]
}
