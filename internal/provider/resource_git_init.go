/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
	"os"
)

type resourceGitInitType struct{}

type resourceGitInit struct {
	p gitProvider
}

type resourceGitInitSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Bare      types.Bool   `tfsdk:"bare"`
}

func (r *resourceGitInitType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Initializes a Git repository similar to 'git init'.",
		MarkdownDescription: "Initializes a Git repository similar to `git init`.",
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
				Description:         "The import ID to import this resource which is equal to the value of the 'directory' attribute.",
				MarkdownDescription: "The import ID to import this resource which is equal to the value of the `directory` attribute.",
				Type:                types.StringType,
				Computed:            true,
			},
			"bare": {
				Description:         "Whether the created Git repository is bare or not. Defaults to 'false'.",
				MarkdownDescription: "Whether the created Git repository is bare or not. Defaults to `false`.",
				Type:                types.BoolType,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.Bool{Value: false}),
					resource.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (r *resourceGitInitType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return &resourceGitInit{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *resourceGitInit) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create resource git_init")

	var inputs resourceGitInitSchema
	var state resourceGitInitSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// NOTE: It seems default values are not working?
	if inputs.Bare.IsNull() {
		inputs.Bare = types.Bool{Value: false}
	}

	directory := inputs.Directory.Value
	bare := inputs.Bare.Value

	_, err := git.PlainInit(directory, bare)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot create repository",
			"Could not create repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created repository", map[string]interface{}{
		"directory": directory,
		"bare":      bare,
	})

	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.Bare = inputs.Bare

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitInit) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read resource git_init")

	var state resourceGitInitSchema
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var newState resourceGitInitSchema
	newState.Directory = state.Directory

	worktree, err := getWorktree(repository, &resp.Diagnostics)
	if err != nil {
		return
	}
	newState.Bare = types.Bool{Value: worktree == nil}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitInit) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update resource git_init")
	updatedUsingPlan(ctx, &req, resp, &resourceGitInitSchema{})
}

func (r *resourceGitInit) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete resource git_init")

	var state resourceGitInitSchema
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.Value
	bare := state.Bare.Value

	if !bare {
		repository := openRepository(ctx, directory, &resp.Diagnostics)

		if repository.Storer != nil {
			storage, ok := repository.Storer.(*filesystem.Storage)

			if ok {
				err := os.RemoveAll(storage.Filesystem().Root())
				if err != nil {
					resp.Diagnostics.AddError(
						"Cannot delete repository",
						"Could not delete git repository ["+directory+"] because of: "+err.Error(),
					)
					return
				}
			}
		}
	}
}

func (r *resourceGitInit) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "ImportState resource git_init")

	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Unexpected import identifier",
			fmt.Sprintf("Expected import identifier with format: 'path/to/git/repository' Got: '%q'", req.ID),
		)
		return
	}

	repository := openRepository(ctx, req.ID, &resp.Diagnostics)
	if repository == nil {
		return
	}

	var state resourceGitInitSchema
	state.Directory = types.String{Value: req.ID}
	state.Id = types.String{Value: req.ID}
	worktree, err := getWorktree(repository, &resp.Diagnostics)
	if err != nil {
		return
	}
	state.Bare = types.Bool{Value: worktree == nil}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
