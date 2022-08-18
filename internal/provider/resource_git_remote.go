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

func (r *resourceGitRemoteType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Manages remotes in a Git repository with `git remote`.",
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
				MarkdownDescription: "The import ID to import this resource which has the form `'directory|name'`",
				Type:                types.StringType,
				Computed:            true,
			},
			"name": {
				Description: "The name of the Git remote to manage.",
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
				Description: "The URLs of the Git remote to manage. The first URL will be a fetch/pull URL. All other URLs will be push only.",
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
	tflog.Debug(ctx, "Create git_remote")

	var inputs resourceGitRemoteSchema

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

	var state resourceGitRemoteSchema
	state.Directory = inputs.Directory
	state.Id = types.String{Value: fmt.Sprintf("%s|%s", directory, name)}
	state.Name = inputs.Name
	state.Urls = inputs.Urls

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitRemote) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read git_remote")

	var state resourceGitRemoteSchema
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.Value
	name := state.Name.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	remote := getRemote(ctx, repository, name, &resp.Diagnostics)
	if remote == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var newState resourceGitRemoteSchema
	newState.Directory = state.Directory
	newState.Id = types.String{Value: fmt.Sprintf("%s|%s", directory, name)}
	newState.Name = state.Name
	newState.Urls = stringsToList(remote.Config().URLs)

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitRemote) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update git_remote")

	var inputs resourceGitRemoteSchema

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

	remote := getRemote(ctx, repository, name, &resp.Diagnostics)
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
	cfg.Remotes[name] = remoteConfig

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

	var state resourceGitRemoteSchema
	state.Directory = inputs.Directory
	state.Id = types.String{Value: fmt.Sprintf("%s|%s", directory, name)}
	state.Name = inputs.Name
	state.Urls = inputs.Urls

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitRemote) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete git_remote")

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
	tflog.Trace(ctx, "deleted remote", map[string]interface{}{
		"directory": directory,
		"remote":    name,
	})
}

func (r *resourceGitRemote) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "ImportState git_remote")

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
	name := idParts[1]
	tflog.Trace(ctx, "parsed import ID", map[string]interface{}{
		"directory": directory,
		"remote":    name,
	})

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	remote := getRemote(ctx, repository, name, &resp.Diagnostics)
	if remote == nil {
		return
	}

	var state resourceGitRemoteSchema
	state.Directory = types.String{Value: directory}
	state.Id = types.String{Value: id}
	state.Name = types.String{Value: name}
	state.Urls = stringsToList(remote.Config().URLs)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
