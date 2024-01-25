// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Cloudbase Solutions SRL
//
//	Licensed under the Apache License, Version 2.0 (the "License"); you may
//	not use this file except in compliance with the License. You may obtain
//	a copy of the License at
//
//	     http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//	WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//	License for the specific language governing permissions and limitations
//	under the License.

package util

import (
	"fmt"
	"strings"

	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm-provider-gcp/config"
	"github.com/cloudbase/garm-provider-gcp/internal/spec"
	"google.golang.org/protobuf/proto"
)

const (
	startup = "user-data"
	netTier = "PREMIUM"
)

func getMachineType(zone, flavor string) string {
	machine := fmt.Sprintf("zones/%s/machineTypes/%s", zone, flavor)
	return machine
}

func GetInstanceName(name string) string {
	lowerName := strings.ToLower(name)
	return lowerName
}

func GenerateInstance(cfg *config.Config, spec *spec.RunnerSpec, udata string) *computepb.Instance {
	name := GetInstanceName(spec.BootstrapParams.Name)
	nicType := "VIRTIO_NET"
	accessConfig := &computepb.AccessConfig{
		NetworkTier: proto.String(netTier),
	}
	Labels := map[string]string{
		"garmpoolid":       spec.BootstrapParams.PoolID,
		"garmcontrollerid": spec.ControllerID,
	}
	inst := &computepb.Instance{
		Name:        proto.String(name),
		MachineType: proto.String(getMachineType(cfg.Zone, spec.BootstrapParams.Flavor)),
		Disks: []*computepb.AttachedDisk{
			{
				Boot: proto.Bool(true),
				InitializeParams: &computepb.AttachedDiskInitializeParams{
					DiskSizeGb:  proto.Int64(spec.DiskSize),
					SourceImage: proto.String(spec.BootstrapParams.Image),
				},
				AutoDelete: proto.Bool(true),
			},
		},
		NetworkInterfaces: []*computepb.NetworkInterface{
			{
				Network:       proto.String(cfg.NetworkID),
				NicType:       proto.String(nicType),
				AccessConfigs: []*computepb.AccessConfig{accessConfig},
				Subnetwork:    proto.String("projects/garm-testing/regions/europe-west1/subnetworks/garm"),
			},
		},
		Metadata: &computepb.Metadata{
			Items: []*computepb.Items{
				{
					Key:   proto.String(startup),
					Value: proto.String(udata),
				},
			},
		},
		Labels: Labels,
		ServiceAccounts: []*computepb.ServiceAccount{
			{
				Email: proto.String("626797023368-compute@developer.gserviceaccount.com"),
				Scopes: []string{
					"https://www.googleapis.com/auth/cloud-platform",
					"https://www.googleapis.com/auth/compute",
					"https://www.googleapis.com/auth/devstorage.read_only",
					"https://www.googleapis.com/auth/logging.write",
					"https://www.googleapis.com/auth/monitoring.write",
				},
			},
		},
	}

	return inst
}

func GcpInstanceToParamsInstance(gcpInstance *computepb.Instance) (params.ProviderInstance, error) {
	if gcpInstance == nil {
		return params.ProviderInstance{}, fmt.Errorf("instance ID is nil")
	}
	details := params.ProviderInstance{
		ProviderID: GetInstanceName(*gcpInstance.Name),
		Name:       GetInstanceName(*gcpInstance.Name),
		OSType:     params.OSType("debian"),
		OSArch:     params.OSArch("amd64"),
		Status:     params.InstanceRunning,
	}
	return details, nil
}
