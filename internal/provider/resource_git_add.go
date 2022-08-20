/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

type resourceGitAddType struct{}

type resourceGitAdd struct {
	p gitProvider
}

type resourceGitAddSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	File      types.String `tfsdk:"file"`
	FileSHA1  types.String `tfsdk:"file_sha1"`
}

func (r *resourceGitAddType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Add file contents to the index using `git add`.",
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
				MarkdownDescription: "The import ID to import this resource which has the form `'directory|file'`",
				Type:                types.StringType,
				Computed:            true,
			},
			"file": {
				Description: "The file to add to the Git index.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"file_sha1": {
				MarkdownDescription: "The SHA1 checksum of the content in `file`.",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (r *resourceGitAddType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return &resourceGitAdd{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *resourceGitAdd) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create git_add")

	var inputs resourceGitAddSchema
	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	name := inputs.File.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	worktree, err := getWorktree(repository, &resp.Diagnostics)
	if err != nil {
		return
	}
	if worktree == nil {
		resp.Diagnostics.AddError(
			"Cannot add file to bare repository",
			"The repository at ["+directory+"] is bare. Create a worktree first in order to add files to it.",
		)
		return
	}

	sha1 := readFileSha1(err, worktree, name, &resp.Diagnostics)
	if sha1 == "" {
		return
	}

	err = addFile(worktree, name, &resp.Diagnostics)
	if err != nil {
		return
	}

	var state resourceGitAddSchema
	state.Directory = inputs.Directory
	state.Id = types.String{Value: fmt.Sprintf("%s|%s", directory, name)}
	state.File = inputs.File
	state.FileSHA1 = types.String{Value: sha1}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitAdd) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read git_add")

	var state resourceGitAddSchema
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.Value
	name := state.File.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	worktree, err := getWorktree(repository, &resp.Diagnostics)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	if worktree == nil {
		resp.Diagnostics.AddError(
			"Cannot add file to bare repository",
			"The repository at ["+directory+"] is bare. Create a worktree first in order to add files to it.",
		)
		resp.State.RemoveResource(ctx)
		return
	}

	sha1 := readFileSha1(err, worktree, name, &resp.Diagnostics)
	if sha1 == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	var newState resourceGitAddSchema
	newState.Directory = state.Directory
	newState.Id = types.String{Value: fmt.Sprintf("%s|%s", directory, name)}
	newState.File = state.File
	newState.FileSHA1 = types.String{Value: sha1}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitAdd) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update git_add")

	var inputs resourceGitAddSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	name := inputs.File.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	worktree, err := getWorktree(repository, &resp.Diagnostics)
	if err != nil {
		return
	}
	if worktree == nil {
		resp.Diagnostics.AddError(
			"Cannot add file to bare repository",
			"The repository at ["+directory+"] is bare. Create a worktree first in order to add files to it.",
		)
		return
	}

	sha1 := readFileSha1(err, worktree, name, &resp.Diagnostics)
	if sha1 == "" {
		return
	}

	err = addFile(worktree, name, &resp.Diagnostics)
	if err != nil {
		return
	}

	var state resourceGitAddSchema
	state.Directory = inputs.Directory
	state.Id = types.String{Value: fmt.Sprintf("%s|%s", directory, name)}
	state.File = inputs.File
	state.FileSHA1 = types.String{Value: sha1}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitAdd) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete git_add")
	// NO-OP: Terraform removes the state automatically for us
}

func (r *resourceGitAdd) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "ImportState git_add")

	id := req.ID
	idParts := strings.Split(id, "|")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected import identifier",
			fmt.Sprintf("Expected import identifier with format: 'path/to/your/git/repository|path/to/your/file' Got: '%q'", id),
		)
		return
	}

	directory := idParts[0]
	name := idParts[1]
	tflog.Trace(ctx, "parsed import ID", map[string]interface{}{
		"directory": directory,
		"file":      name,
	})

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	worktree, err := getWorktree(repository, &resp.Diagnostics)
	if err != nil {
		return
	}
	if worktree == nil {
		resp.Diagnostics.AddError(
			"Cannot add file to bare repository",
			"The repository at ["+directory+"] is bare. Create a worktree first in order to add files to it.",
		)
		return
	}

	sha1 := readFileSha1(err, worktree, name, &resp.Diagnostics)
	if sha1 == "" {
		return
	}

	var state resourceGitAddSchema
	state.Directory = types.String{Value: directory}
	state.Id = types.String{Value: id}
	state.File = types.String{Value: name}
	state.FileSHA1 = types.String{Value: sha1}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
