/*
Copyright (c) 2015 VMware, Inc. All Rights Reserved.

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

package object

import (
	"fmt"

	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
)

// ResourcePool represents a client to manage a resource pool
type ResourcePool struct {
	Common

	InventoryPath string
}

// String has a string representation of the resource pool
func (p ResourcePool) String() string {
	if p.InventoryPath == "" {
		return p.Common.String()
	}
	return fmt.Sprintf("%v @ %v", p.Common, p.InventoryPath)
}

// NewResourcePool creates a new resource pool client
func NewResourcePool(c *vim25.Client, ref types.ManagedObjectReference) *ResourcePool {
	return &ResourcePool{
		Common: NewCommon(c, ref),
	}
}

// Name returns the name of the resource poool
func (p ResourcePool) Name(ctx context.Context) (string, error) {
	var o mo.ResourcePool

	err := p.Properties(ctx, p.Reference(), []string{"name"}, &o)
	if err != nil {
		return "", err
	}

	return o.Name, nil
}

// ImportVApp imports a vApp into the resource pool
func (p ResourcePool) ImportVApp(ctx context.Context, spec types.BaseImportSpec, folder *Folder, host *HostSystem) (*HTTPNfcLease, error) {
	req := types.ImportVApp{
		This: p.Reference(),
		Spec: spec,
	}

	if folder != nil {
		ref := folder.Reference()
		req.Folder = &ref
	}

	if host != nil {
		ref := host.Reference()
		req.Host = &ref
	}

	res, err := methods.ImportVApp(ctx, p.c, &req)
	if err != nil {
		return nil, err
	}

	return NewHTTPNfcLease(p.c, res.Returnval), nil
}

// Create a resource pool
func (p ResourcePool) Create(ctx context.Context, name string, spec types.ResourceConfigSpec) (*ResourcePool, error) {
	req := types.CreateResourcePool{
		This: p.Reference(),
		Name: name,
		Spec: spec,
	}

	res, err := methods.CreateResourcePool(ctx, p.c, &req)
	if err != nil {
		return nil, err
	}

	return NewResourcePool(p.c, res.Returnval), nil
}

// CreateVApp in this resource pool
func (p ResourcePool) CreateVApp(ctx context.Context, name string, resSpec types.ResourceConfigSpec, configSpec types.VAppConfigSpec, folder *Folder) (*VirtualApp, error) {
	req := types.CreateVApp{
		This:       p.Reference(),
		Name:       name,
		ResSpec:    resSpec,
		ConfigSpec: configSpec,
	}

	if folder != nil {
		ref := folder.Reference()
		req.VmFolder = &ref
	}

	res, err := methods.CreateVApp(ctx, p.c, &req)
	if err != nil {
		return nil, err
	}

	return NewVirtualApp(p.c, res.Returnval), nil
}

// UpdateConfig updates the config for this resource pool
func (p ResourcePool) UpdateConfig(ctx context.Context, name string, config *types.ResourceConfigSpec) error {
	req := types.UpdateConfig{
		This:   p.Reference(),
		Name:   name,
		Config: config,
	}

	if config != nil && config.Entity == nil {
		ref := p.Reference()

		// Create copy of config so changes won't leak back to the caller
		newConfig := *config
		newConfig.Entity = &ref
		req.Config = &newConfig
	}

	_, err := methods.UpdateConfig(ctx, p.c, &req)
	return err
}

// DestroyChildren destroys the children for this resource pool
func (p ResourcePool) DestroyChildren(ctx context.Context) error {
	req := types.DestroyChildren{
		This: p.Reference(),
	}

	_, err := methods.DestroyChildren(ctx, p.c, &req)
	return err
}

// Destroy this resource pool
func (p ResourcePool) Destroy(ctx context.Context) (*Task, error) {
	req := types.Destroy_Task{
		This: p.Reference(),
	}

	res, err := methods.Destroy_Task(ctx, p.c, &req)
	if err != nil {
		return nil, err
	}

	return NewTask(p.c, res.Returnval), nil
}
