/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5/config"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	Urls      types.List   `tfsdk:"urls"`
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
					resource.RequiresReplace(),
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
					resource.RequiresReplace(),
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

func (r *resourceGitRemoteType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return &resourceGitRemote{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *resourceGitRemote) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	tflog.Trace(ctx, "opened repository", map[string]interface{}{
		"directory": directory,
	})

	urls := make([]string, len(inputs.Urls.Elems))
	resp.Diagnostics.Append(inputs.Urls.ElementsAs(ctx, &urls, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := repository.CreateRemote(&config.RemoteConfig{
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
		"remote":    name,
	})

	output.Directory = inputs.Directory
	output.Id = inputs.Name
	output.Name = inputs.Name
	output.Urls = inputs.Urls

	diags = resp.State.Set(ctx, &output)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitRemote) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Git remote")

	var state resourceGitRemoteSchema
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.Value
	remoteName := state.Name.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	remote := getRemote(ctx, repository, remoteName, &resp.Diagnostics)
	if remote == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var newState resourceGitRemoteSchema
	newState.Directory = state.Directory
	newState.Id = state.Name
	newState.Name = state.Name
	newState.Urls = stringsToList(remote.Config().URLs)

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitRemote) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	tflog.Trace(ctx, "read repository config", map[string]interface{}{
		"directory": directory,
	})

	urls := make([]string, len(inputs.Urls.Elems))
	resp.Diagnostics.Append(inputs.Urls.ElementsAs(ctx, &urls, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	remoteConfig := remote.Config()
	remoteConfig.URLs = urls
	cfg.Remotes[remoteName] = remoteConfig

	err = repository.SetConfig(cfg)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot write repository config",
			"Could not write repository config of ["+directory+"] because of: "+err.Error(),
		)
		return
	}
	tflog.Trace(ctx, "wrote repository config", map[string]interface{}{
		"directory": directory,
	})

	state.Directory = inputs.Directory
	state.Id = inputs.Name
	state.Name = inputs.Name
	state.Urls = inputs.Urls

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitRemote) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Removing Git remote")

	var state resourceGitRemoteSchema

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.Value
	remoteName := state.Name.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	err := repository.DeleteRemote(remoteName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot delete remote",
			"Could not delete remote ["+remoteName+"] in git repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}
	tflog.Trace(ctx, "deleted remote", map[string]interface{}{
		"directory": directory,
		"remote":    remoteName,
	})
}

func (r *resourceGitRemote) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	idParts := strings.Split(id, "|")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected import identifier",
			fmt.Sprintf("Expected import identifier with format: 'path/to/your/git/repository|name-of-your-remote' Got: '%q'", id),
		)
		return
	}

	directory := idParts[0]
	remoteName := idParts[1]
	tflog.Trace(ctx, "parsed import ID", map[string]interface{}{
		"directory": directory,
		"remote":    remoteName,
	})

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	remote := getRemote(ctx, repository, remoteName, &resp.Diagnostics)
	if remote == nil {
		return
	}

	var state resourceGitRemoteSchema
	state.Directory = types.String{Value: directory}
	state.Id = types.String{Value: remoteName}
	state.Name = types.String{Value: remoteName}
	state.Urls = stringsToList(remote.Config().URLs)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
