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
)

func getMachineType(zone, flavor string) *string {
	machine := fmt.Sprintf("zones/%s/machineTypes/%s", zone, flavor)
	return &machine
}

func GetInstanceName(name string) string {
	lowerName := strings.ToLower(name)
	return lowerName
}

func GenerateInstance(cfg *config.Config, spec *spec.RunnerSpec, udata string) *computepb.Instance {
	boot := true
	name := GetInstanceName(spec.BootstrapParams.Name)
	startup := "startup-script"
	inst := &computepb.Instance{
		Name:        &name,
		MachineType: getMachineType(cfg.Zone, spec.BootstrapParams.Flavor),
		Disks: []*computepb.AttachedDisk{
			{
				Boot: &boot,
				InitializeParams: &computepb.AttachedDiskInitializeParams{
					SourceImage: &spec.BootstrapParams.Image,
				},
			},
		},
		NetworkInterfaces: []*computepb.NetworkInterface{
			{
				Network: &spec.NetworkID,
			},
		},
		Metadata: &computepb.Metadata{
			Items: []*computepb.Items{
				{
					Key:   &startup,
					Value: &udata,
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
	disks := gcpInstance.GetDisks()
	disk := disks[0]
	details := params.ProviderInstance{
		ProviderID: *gcpInstance.Name,
		Name:       *gcpInstance.Name,
		OSType:     params.OSType(disk.GuestOsFeatures[0].GetType()),
		OSArch:     params.OSArch(*gcpInstance.CpuPlatform),
		//*gcpInstance.ResourceStatus.String()
		Status: params.InstanceRunning,
	}
	return details, nil
}
