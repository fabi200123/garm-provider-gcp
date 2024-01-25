// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Cloudbase Solutions SRL
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package spec

import (
	"encoding/json"
	"fmt"

	"github.com/cloudbase/garm-provider-common/cloudconfig"
	"github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm-provider-common/util"
	"github.com/cloudbase/garm-provider-gcp/config"
)

const (
	defaultDiskSizeGB int64 = 100
)

func newExtraSpecsFromBootstrapData(data params.BootstrapInstance) (*extraSpecs, error) {
	spec := &extraSpecs{}

	if len(data.ExtraSpecs) > 0 {
		if err := json.Unmarshal(data.ExtraSpecs, spec); err != nil {
			return nil, fmt.Errorf("failed to unmarshal extra specs: %w", err)
		}
	}

	return spec, nil
}

type extraSpecs struct {
	DiskSize  int64   `json:"disksize,omitempty"`
	NetworkID *string `json:"network_id,omitempty"`
}

func GetRunnerSpecFromBootstrapParams(cfg *config.Config, data params.BootstrapInstance, controllerID string) (*RunnerSpec, error) {
	tools, err := util.GetTools(data.OSType, data.OSArch, data.Tools)
	if err != nil {
		return nil, fmt.Errorf("failed to get tools: %s", err)
	}

	extraSpecs, err := newExtraSpecsFromBootstrapData(data)
	if err != nil {
		return nil, fmt.Errorf("error loading extra specs: %w", err)
	}

	spec := &RunnerSpec{
		Zone:            cfg.Zone,
		Tools:           tools,
		BootstrapParams: data,
		NetworkID:       cfg.NetworkID,
		ControllerID:    controllerID,
		DiskSize:        defaultDiskSizeGB,
	}

	spec.MergeExtraSpecs(extraSpecs)

	return spec, nil
}

type RunnerSpec struct {
	Zone            string
	Tools           params.RunnerApplicationDownload
	BootstrapParams params.BootstrapInstance
	NetworkID       string
	ControllerID    string
	DiskSize        int64
}

func (r *RunnerSpec) MergeExtraSpecs(extraSpecs *extraSpecs) {
	if extraSpecs.NetworkID != nil && *extraSpecs.NetworkID != "" {
		r.NetworkID = *extraSpecs.NetworkID
	}
	if extraSpecs.DiskSize > 0 {
		r.DiskSize = extraSpecs.DiskSize
	}
}

func (r *RunnerSpec) Validate() error {
	if r.Zone == "" {
		return fmt.Errorf("missing zone")
	}
	if r.NetworkID == "" {
		return fmt.Errorf("missing network id")
	}
	if r.ControllerID == "" {
		return fmt.Errorf("missing controller id")
	}

	return nil
}

func (r RunnerSpec) ComposeUserData() (string, error) {
	switch r.BootstrapParams.OSType {
	case params.Linux, params.Windows:
		udata, err := cloudconfig.GetCloudConfig(r.BootstrapParams, r.Tools, r.BootstrapParams.Name)
		if err != nil {
			return "", fmt.Errorf("failed to generate userdata: %w", err)
		}
		return udata, nil
	}
	return "", fmt.Errorf("unsupported OS type for cloud config: %s", r.BootstrapParams.OSType)
}
