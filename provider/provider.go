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

package provider

import (
	"context"
	"fmt"

	"github.com/cloudbase/garm-provider-common/execution"
	"github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm-provider-gcp/internal/client"
)

var _ execution.ExternalProvider = &GceProvider{}

func NewGceProvider(cfgFile string, controllerID string) (*GceProvider, error) {
	return &GceProvider{}, nil
}

type GceProvider struct {
	gcpCli *client.GcpCli
}

func (g *GceProvider) CreateInstance(ctx context.Context, bootstrapParams params.BootstrapInstance) (params.ProviderInstance, error) {
	return params.ProviderInstance{}, nil
}

func (g *GceProvider) GetInstance(ctx context.Context, instance string) (params.ProviderInstance, error) {
	return params.ProviderInstance{}, nil
}

func (g *GceProvider) DeleteInstance(ctx context.Context, instance string) error {
	err := g.gcpCli.DeleteInstance(ctx, instance)
	if err != nil {
		return fmt.Errorf("error deleting instance: %w", err)
	}
	return nil
}

func (g *GceProvider) ListInstances(ctx context.Context, poolID string) ([]params.ProviderInstance, error) {
	return []params.ProviderInstance{}, nil
}

func (g *GceProvider) RemoveAllInstances(ctx context.Context) error {
	return nil
}

func (g *GceProvider) Stop(ctx context.Context, instance string, force bool) error {
	return nil

}

func (g *GceProvider) Start(ctx context.Context, instance string) error {
	return nil
}
