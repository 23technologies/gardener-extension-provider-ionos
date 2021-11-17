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
	"encoding/base64"
	"errors"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis"
	ionosapiwrapper "github.com/23technologies/ionos-api-wrapper/pkg"
	ionossdk "github.com/ionos-cloud/sdk-go/v5"
)

// Constant ionosPasswordGeneratedLength is the length of the generated random password
const ionosPasswordGeneratedLength = 32
// Constant ionosVolumeType is the volume type
const ionosVolumeType = "SSD"

func createDHCPServer(ctx context.Context, client *ionossdk.APIClient, datacenterID, zone string, configuration *apis.DHCPServerConfiguration, networkID string) (string, error) {
	imageID, err := apis.FindMachineImageName(ctx, client, zone, configuration.Image.Name, configuration.Image.Version, "")
	if nil != err {
		return "", err
	}

	image, _, err := client.ImageApi.ImagesFindById(ctx, imageID).Depth(1).Execute()
	if nil != err {
		return "", err
	} else if (!image.Properties.HasCloudInit() || "NONE" == *image.Properties.CloudInit) {
		return "", errors.New("imageID given doesn't belong to a cloud-init enabled image")
	}

	password := configuration.Password
	userDataBase64Encoded := base64.StdEncoding.EncodeToString([]byte(configuration.UserData))
	volumeName := "dhcp-server-root-volume"
	volumeSize := configuration.VolumeSize
	volumeType := ionosVolumeType

	if "" == password {
		var passwordBuilder strings.Builder
		seededRand := rand.New(rand.NewSource(time.Now().Unix()))

		for i := 0; i < ionosPasswordGeneratedLength; i++ {
			passwordBuilder.WriteString(string(32 + seededRand.Intn(94)))
		}

		password = passwordBuilder.String()
	}

	if 0 == volumeSize {
		volumeSize = *image.Properties.Size
	} else {
		volumeSize = float32(math.Max(math.Ceil(float64(volumeSize) / 1048576), float64(*image.Properties.Size)))
	}

	volumeProperties := ionossdk.VolumeProperties{
		Type: &volumeType,
		Name: &volumeName,
		Size: &volumeSize,
		Image: &imageID,
		ImagePassword: &password,
		UserData: &userDataBase64Encoded,
	}

	volumeApiCreateRequest := client.VolumeApi.DatacentersVolumesPost(ctx, datacenterID).Depth(0)
	volume, _, err := volumeApiCreateRequest.Volume(ionossdk.Volume{Properties: &volumeProperties}).Execute()
	if nil != err {
		return "", err
	}

	volumeID := *volume.Id

	volume, err = ionosapiwrapper.WaitForVolumeModificationsAndGetResult(ctx, client, datacenterID, volumeID)
	if nil != err {
		return "", err
	}

	cores := int32(configuration.Cores)
	memory := int32(configuration.Memory)

	serverEntities := ionossdk.ServerEntities{
		Volumes: &ionossdk.AttachedVolumes{Items: &[]ionossdk.Volume{ionossdk.Volume{Id: &volumeID}}},
	}

	serverName := "dhcp-server"

	serverProperties := ionossdk.ServerProperties{
		Name:       &serverName,
		Cores:      &cores,
		Ram:        &memory,
		BootVolume: &ionossdk.ResourceReference{Id: &volumeID},
	}

	serverApiCreateRequest := client.ServerApi.DatacentersServersPost(ctx, datacenterID).Depth(0)
	server, _, err := serverApiCreateRequest.Server(ionossdk.Server{Entities: &serverEntities, Properties: &serverProperties}).Execute()
	if nil != err {
		return "", err
	}

	serverID := *server.Id

	err = ionosapiwrapper.WaitForServerModifications(ctx, client, datacenterID, serverID)
	if nil != err {
		return "", err
	}

	_, _, err = client.ServerApi.DatacentersServersStopPost(ctx, datacenterID, serverID).Execute()
	if nil != err {
		return "", err
	}

	err = ionosapiwrapper.WaitForServerModifications(ctx, client, datacenterID, serverID)
	if nil != err {
		return "", err
	}

	err = ionosapiwrapper.AttachLANToServerWithoutDHCP(ctx, client, datacenterID, serverID, networkID)
	if nil != err {
		return "", err
	}

	_, _, err = client.ServerApi.DatacentersServersStartPost(ctx, datacenterID, serverID).Execute()
	if nil != err {
		return "", err
	}

	err = ionosapiwrapper.WaitForServerModifications(ctx, client, datacenterID, serverID)
	if nil != err {
		return "", err
	}

	return serverID, nil
}

// deleteServerID
//
// PARAMETERS
// ctx       context.Context Execution context
// client    *hcloud.Client  HCloud client
// namespace string          Shoot namespace
// networks  *apis.Networks  Networks struct
func deleteServerID(ctx context.Context, client *ionossdk.APIClient, datacenterID, serverID string) error {
	_, httpResponse, err := client.ServerApi.DatacentersServersStopPost(ctx, datacenterID, serverID).Execute()
	if nil != err {
		if 404 == httpResponse.StatusCode {
			return nil
		} else {
			return err
		}
	}

	err = ionosapiwrapper.WaitForServerModifications(ctx, client, datacenterID, serverID)
	if nil != err {
		return err
	}

	server, _, err := client.ServerApi.DatacentersServersFindById(ctx, datacenterID, serverID).Depth(3).Execute()
	if nil != err {
		return err
	}

	for _, volume := range *server.Entities.Volumes.Items {
		_, _, err := client.VolumeApi.DatacentersVolumesDelete(ctx, datacenterID, *volume.Id).Depth(0).Execute()
		if nil != err {
			return err
		}
	}

	_, _, err = client.ServerApi.DatacentersServersDelete(ctx, datacenterID, serverID).Depth(0).Execute()
	if nil != err {
		return err
	}

	return nil
}

func EnsureDHCPServer(ctx context.Context, client *ionossdk.APIClient, datacenterID, zone string, infraConfig *apis.InfrastructureConfig, infraStatus *apis.InfrastructureStatus, networkID string) (string, error) {
	cidr := infraConfig.Networks.Workers
	configuration := infraConfig.DHCPServerConfiguration
	var configurationStatus *apis.DHCPServerConfigurationStatus

	if nil != infraStatus {
		configurationStatus = infraStatus.DHCPServerConfiguration
	}

	serverID := ""
	oldServerID := ""

	if nil != configurationStatus {
		var err error

		if "" == configurationStatus.ServerID {
			configurationStatus = nil
		} else {
			_, _, err = client.ServerApi.DatacentersServersFindById(ctx, datacenterID, configurationStatus.ServerID).Depth(0).Execute()

			if nil != err || cidr != configurationStatus.Cidr {
				oldServerID = configurationStatus.ServerID
			}
		}
	}

	if nil == configurationStatus || "" != oldServerID {
		if "" != oldServerID {
			err := deleteServerID(ctx, client, datacenterID, oldServerID)
			if nil != err {
				return "", err
			}
		}

		if "" != networkID {
			newServerID, err := createDHCPServer(ctx, client, datacenterID, zone, configuration, networkID)
			if nil != err {
				return "", err
			}

			serverID = newServerID
		}
	}

	return serverID, nil
}

func EnsureDHCPServerDeleted(ctx context.Context, client *ionossdk.APIClient, datacenterID, serverID string) error {
	return deleteServerID(ctx, client, datacenterID, serverID)
}
