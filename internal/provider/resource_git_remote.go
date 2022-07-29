/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5/config"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

type resourceGitRemoteType struct{}

type resourceGitRemote struct {
	p gitProvider
}

type resourceGitRemoteSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Urls      []string     `tfsdk:"urls"`
}

func (c *resourceGitRemoteType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Adds a new Git remote to a repository.",
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description: "The path to the local Git repository.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"id": {
				MarkdownDescription: "`DEPRECATED`: Only added in order to use the sdkv2 test framework. The path to the local Git repository.",
				Type:                types.StringType,
				Computed:            true,
			},
			"name": {
				Description: "The name of the Git remote to add.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"urls": {
				Description: "The URLs of the Git remote to add.",
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Required: true,
			},
		},
	}, nil
}

func (r *resourceGitRemoteType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return &resourceGitRemote{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *resourceGitRemote) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, "Creating Git remote")

	var inputs resourceGitRemoteSchema
	var output resourceGitRemoteSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	name := inputs.Name.Value
	urls := inputs.Urls

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	tflog.Trace(ctx, "opened repository", map[string]interface{}{
		"directory": directory,
	})

	remote, err := repository.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: urls,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot create remote",
			"Could not create remote ["+name+"] in git repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created remote", map[string]interface{}{
		"directory": directory,
		"remote":    remote.Config().Name,
	})

	output.Directory = types.String{Value: directory}
	output.Id = types.String{Value: remote.Config().Name}
	output.Name = types.String{Value: remote.Config().Name}
	output.Urls = urls

	diags = resp.State.Set(ctx, &output)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitRemote) Read(ctx context.Context, _ tfsdk.ReadResourceRequest, _ *tfsdk.ReadResourceResponse) {
	tflog.Debug(ctx, "Reading Git remote from state")
	// NO-OP: all there is to read is in the State, and response is already populated with that.
}

func (r *resourceGitRemote) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "Updating Git remote")

	var inputs resourceGitRemoteSchema
	var state resourceGitRemoteSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	remoteName := inputs.Name.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	remote := getRemote(ctx, repository, remoteName, &resp.Diagnostics)
	if remote == nil {
		return
	}

	cfg, err := repository.Config()
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read repository config",
			"Could not read repository config of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	remoteConfig := remote.Config()
	remoteConfig.URLs = inputs.Urls
	cfg.Remotes[remoteName] = remoteConfig

	err = repository.SetConfig(cfg)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot write repository config",
			"Could not write repository config of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	state.Directory = types.String{Value: directory}
	state.Id = types.String{Value: remoteName}
	state.Name = types.String{Value: remoteName}
	state.Urls = inputs.Urls

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitRemote) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "Removing Git remote")

	var state resourceGitRemoteSchema

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.Value
	name := state.Name.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	err := repository.DeleteRemote(name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot delete remote",
			"Could not delete remote ["+name+"] in git repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}
}

func (r *resourceGitRemote) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	id := req.ID
	idParts := strings.Split(id, "|")

	if len(idParts) < 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected import identifier",
			fmt.Sprintf("Expected import identifier with format: 'directory|remote-name|first-url,second-url,...' Got: %q", id),
		)
		return
	}

	var state resourceGitRemoteSchema

	state.Directory = types.String{Value: idParts[0]}
	state.Id = types.String{Value: idParts[1]}
	state.Name = types.String{Value: idParts[1]}
	state.Urls = strings.Split(idParts[2], ",")

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
