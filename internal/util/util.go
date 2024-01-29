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
	startup          string = "startup-script"
	netTier          string = "PREMIUM"
	accessConfigType string = "ONE_TO_ONE_NAT"
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
				Network: proto.String(cfg.NetworkID),
				NicType: proto.String(spec.NicType),
				AccessConfigs: []*computepb.AccessConfig{
					{
						// The type of configuration. In accessConfigs (IPv4), the default and only option is ONE_TO_ONE_NAT.
						Type:        proto.String(accessConfigType),
						NetworkTier: proto.String(netTier),
					},
				},
				Subnetwork: &spec.SubnetworkID,
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
		Labels: map[string]string{
			"garmpoolid":       spec.BootstrapParams.PoolID,
			"garmcontrollerid": spec.ControllerID,
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
		OSType:     params.OSType("linux"),
		OSArch:     params.OSArch("amd64"),
	}

	switch gcpInstance.GetStatus() {
	case "RUNNING":
		details.Status = params.InstanceRunning
	case "STOPPING", "TERMINATED", "SUSPENDED":
		details.Status = params.InstanceStopped
	case "PROVISIONING":
		details.Status = params.InstancePendingCreate
	case "STAGING":
		details.Status = params.InstanceCreating
	default:
		details.Status = params.InstanceStatusUnknown
	}

	return details, nil
}
