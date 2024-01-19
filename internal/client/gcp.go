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
	"os"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/cloudbase/garm-provider-gcp/config"
	"github.com/cloudbase/garm-provider-gcp/internal/spec"
	"github.com/cloudbase/garm-provider-gcp/internal/util"
	"golang.org/x/oauth2/google"
	gcompute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func NewGcpCli(ctx context.Context, cfg *config.Config) (*GcpCli, error) {
	jsonKey, err := os.ReadFile(cfg.CredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON key file: %w", err)
	}
	config, err := google.JWTConfigFromJSON(jsonKey, gcompute.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT config: %w", err)
	}
	// Create an HTTP client using the JWT Config
	client := config.Client(ctx)

	// Now use this client to create a Compute Engine client
	computeClient, err := compute.NewInstancesRESTClient(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creating compute service: %w", err)
	}
	gcpCli := &GcpCli{
		cfg:    cfg,
		client: computeClient,
	}

	return gcpCli, nil
}

type GcpCli struct {
	cfg *config.Config

	client *compute.InstancesClient
}

func (g *GcpCli) CreateInstance(ctx context.Context, spec *spec.RunnerSpec) (string, error) {
	if spec == nil {
		return "", fmt.Errorf("invalid nil runner spec")
	}

	udata, err := spec.ComposeUserData()
	if err != nil {
		return "", fmt.Errorf("failed to compose user data: %w", err)
	}

	inst := util.GenerateInstance(g.cfg, spec, udata)

	insertReq := &computepb.InsertInstanceRequest{
		Project:          g.cfg.ProjectId,
		Zone:             g.cfg.Zone,
		InstanceResource: inst,
	}

	op, err := g.client.Insert(ctx, insertReq)
	if err != nil {
		return "", fmt.Errorf("failed to create instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return "", fmt.Errorf("failed to wait for operation: %w", err)
	}

	return spec.BootstrapParams.Name, nil
}

func (g *GcpCli) GetInstance(ctx context.Context, instanceName string) (*computepb.Instance, error) {
	req := &computepb.GetInstanceRequest{
		Project:  g.cfg.ProjectId,
		Zone:     g.cfg.Zone,
		Instance: instanceName,
	}

	instance, err := g.client.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %v", err)
	}

	return instance, nil
}

func (g *GcpCli) DeleteInstance(ctx context.Context, instance string) error {
	req := &computepb.DeleteInstanceRequest{
		Instance: util.GetInstanceName(instance),
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

func (g *GcpCli) ListInstances(ctx context.Context, poolID string) ([]string, error) {
	return nil, nil
}

func (g *GcpCli) StopInstance(ctx context.Context, instance string) error {
	req := &computepb.StopInstanceRequest{
		Instance: util.GetInstanceName(instance),
		Project:  g.cfg.ProjectId,
		Zone:     g.cfg.Zone,
	}

	op, err := g.client.Stop(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to stop instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	return nil
}

func (g *GcpCli) StartInstance(ctx context.Context, instance string) error {
	req := &computepb.StartInstanceRequest{
		Instance: util.GetInstanceName(instance),
		Project:  g.cfg.ProjectId,
		Zone:     g.cfg.Zone,
	}

	op, err := g.client.Start(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to start instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	return nil
}
