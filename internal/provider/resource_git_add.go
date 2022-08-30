/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
)

type resourceGitAddType struct{}

type resourceGitAdd struct {
	p gitProvider
}

type resourceGitAddSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	All       types.Bool   `tfsdk:"all"`
	ExactPath types.String `tfsdk:"exact_path"`
	GlobPath  types.String `tfsdk:"glob_path"`
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
				MarkdownDescription: "The same value as the `directory` attribute.",
				Type:                types.StringType,
				Computed:            true,
			},
			"all": {
				MarkdownDescription: "Update the index not only where the working tree has a file matching `exact_path` or `glob_path` but also where the index already has an entry. This adds, modifies, and removes index entries to match the working tree. If no paths are given, all files in the entire working tree are updated. Defaults to `true`.",
				Type:                types.BoolType,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.Bool{Value: true}),
					resource.RequiresReplace(),
				},
			},
			"exact_path": {
				Description: "The exact filepath to the file or directory to be added. Conflicts with `glob_path`.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.ConflictsWith(path.MatchRoot("glob_path")),
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"glob_path": {
				MarkdownDescription: "The glob pattern of files or directories to be added. Conflicts with `exact_path`.",
				Type:                types.StringType,
				Computed:            true,
				Optional:            true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.ConflictsWith(path.MatchRoot("exact_path")),
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
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

	// NOTE: It seems default values are not working?
	if inputs.All.IsNull() {
		inputs.All = types.Bool{Value: true}
	}

	options := &git.AddOptions{
		All: inputs.All.Value,
	}
	if !inputs.ExactPath.IsNull() {
		options.Path = inputs.ExactPath.Value
	} else if !inputs.GlobPath.IsNull() {
		options.Glob = inputs.GlobPath.Value
	}

	err = addPaths(worktree, options, &resp.Diagnostics)
	if err != nil {
		return
	}

	var state resourceGitAddSchema
	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.All = inputs.All
	state.ExactPath = inputs.ExactPath
	state.GlobPath = inputs.GlobPath

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitAdd) Read(ctx context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	tflog.Debug(ctx, "Read git_add")
	// NO-OP: All data is already in Terraform state
}

func (r *resourceGitAdd) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update git_add")
	// NO-OP: All attributes require replacement, thus delete/create will be called
}

func (r *resourceGitAdd) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete git_add")
	// NO-OP: Terraform removes the state automatically for us
}
