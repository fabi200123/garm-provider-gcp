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

package client

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/cloudbase/garm-provider-gcp/config"
)

func NewGcpCli(ctx context.Context, cfg *config.Config) (*GcpCli, error) {
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating compute service: %w", err)
	}
	gcpCli := &GcpCli{
		cfg:    cfg,
		client: client,
		zone:   cfg.Zone,
	}

	return gcpCli, nil
}

type GcpCli struct {
	cfg *config.Config

	client *compute.InstancesClient
	zone   string
}

func (g *GcpCli) DeleteInstance(ctx context.Context, instance string) error {
	req := &computepb.DeleteInstanceRequest{
		Instance: instance,
		Project:  g.cfg.ProjectId,
		Zone:     g.cfg.Zone,
	}

	op, err := g.client.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	return nil
}
