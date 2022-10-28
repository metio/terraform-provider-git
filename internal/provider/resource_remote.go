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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

type RemoteResource struct{}

var (
	_ resource.Resource                = (*RemoteResource)(nil)
	_ resource.ResourceWithImportState = (*RemoteResource)(nil)
)

type remoteResourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Urls      types.List   `tfsdk:"urls"`
}

func NewRemoteResource() resource.Resource {
	return &RemoteResource{}
}

func (r *RemoteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote"
}

func (r *RemoteResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Manages remotes in a Git repository similar to 'git remote'.",
		MarkdownDescription: "Manages remotes in a Git repository similar to `git remote`.",
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description:         "The path to the local Git repository.",
				MarkdownDescription: "The path to the local Git repository.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"id": {
				Description:         "The import ID to import this resource which has the form 'directory|name'",
				MarkdownDescription: "The import ID to import this resource which has the form `'directory|name'`",
				Type:                types.StringType,
				Computed:            true,
			},
			"name": {
				Description:         "The name of the Git remote to manage.",
				MarkdownDescription: "The name of the Git remote to manage.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"urls": {
				Description:         "The URLs of the Git remote to manage. The first URL will be a fetch/pull URL. All other URLs will be push only.",
				MarkdownDescription: "The URLs of the Git remote to manage. The first URL will be a fetch/pull URL. All other URLs will be push only.",
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Required: true,
			},
		},
	}, nil
}

func (r *RemoteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create resource git_remote")

	var inputs remoteResourceModel

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()
	name := inputs.Name.ValueString()

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	urls := make([]string, len(inputs.Urls.Elements()))
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

	var state remoteResourceModel
	state.Directory = inputs.Directory
	state.Id = types.StringValue(fmt.Sprintf("%s|%s", directory, name))
	state.Name = inputs.Name
	state.Urls = inputs.Urls

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RemoteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read resource git_remote")

	var state remoteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.ValueString()
	name := state.Name.ValueString()

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

	var newState remoteResourceModel
	newState.Directory = state.Directory
	newState.Id = types.StringValue(fmt.Sprintf("%s|%s", directory, name))
	newState.Name = state.Name
	newState.Urls = StringsToList(remote.Config().URLs)

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RemoteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update resource git_remote")

	var inputs remoteResourceModel

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()
	name := inputs.Name.ValueString()

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

	urls := make([]string, len(inputs.Urls.Elements()))
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

	var state remoteResourceModel
	state.Directory = inputs.Directory
	state.Id = types.StringValue(fmt.Sprintf("%s|%s", directory, name))
	state.Name = inputs.Name
	state.Urls = inputs.Urls

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RemoteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete resource git_remote")

	var state remoteResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.ValueString()
	name := state.Name.ValueString()

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

func (r *RemoteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "ImportState resource git_remote")

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

	var state remoteResourceModel
	state.Directory = types.StringValue(directory)
	state.Id = types.StringValue(id)
	state.Name = types.StringValue(name)
	state.Urls = StringsToList(remote.Config().URLs)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
