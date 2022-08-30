/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
)

type resourceGitCommitType struct{}

type resourceGitCommit struct {
	p gitProvider
}

type resourceGitCommitSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Message   types.String `tfsdk:"message"`
	All       types.Bool   `tfsdk:"all"`
	Author    types.Object `tfsdk:"author"`
	Committer types.Object `tfsdk:"committer"`
	SHA1      types.String `tfsdk:"sha1"`
}

func (r *resourceGitCommitType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Record changes to the repository with `git commit`",
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
			"message": {
				Description: "The commit message to use.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"all": {
				MarkdownDescription: "Automatically stage files that have been modified and deleted, but new files you have not told Git about are not affected. Defaults to `false`.",
				Type:                types.BoolType,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.Bool{Value: false}),
					resource.RequiresReplace(),
				},
			},
			"author": {
				Description: "The original author of the commit. If none is specified, the author will be read from the Git configuration.",
				Computed:    true,
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Description: "The name of the author.",
						Type:        types.StringType,
						Computed:    true,
						Optional:    true,
					},
					"email": {
						Description: "The email address of the author.",
						Type:        types.StringType,
						Computed:    true,
						Optional:    true,
					},
				}),
			},
			"committer": {
				Description: "The person performing the commit. If none is specified, the author is used as committer.",
				Computed:    true,
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Description: "The name of the committer.",
						Type:        types.StringType,
						Computed:    true,
						Optional:    true,
					},
					"email": {
						Description: "The email address of the committer.",
						Type:        types.StringType,
						Computed:    true,
						Optional:    true,
					},
				}),
			},
			"sha1": {
				Description: "The SHA1 hash of the created commit.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func (r *resourceGitCommitType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return &resourceGitCommit{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *resourceGitCommit) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create git_commit")

	var inputs resourceGitCommitSchema
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
			"Cannot commit in bare repository",
			"The repository at ["+directory+"] is bare. Create a worktree first in order to commit changes.",
		)
		return
	}

	// NOTE: It seems default values are not working?
	if inputs.All.IsNull() {
		inputs.All = types.Bool{Value: false}
	}

	status, err := worktree.Status()
	if err != nil {
		return
	}

	var state resourceGitCommitSchema
	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.All = inputs.All
	state.Message = inputs.Message
	state.Author = inputs.Author
	state.Committer = inputs.Committer

	if !status.IsClean() {
		options := createCommitOptions(ctx, inputs)

		hash := createCommit(worktree, inputs.Message.Value, options, &resp.Diagnostics)
		if hash == nil {
			return
		}

		commitObject := getCommit(ctx, repository, hash, &resp.Diagnostics)
		if commitObject == nil {
			return
		}

		state.Author = signatureToObjectWithoutTimestamp(&commitObject.Author)
		state.Committer = signatureToObjectWithoutTimestamp(&commitObject.Committer)
		state.SHA1 = types.String{Value: hash.String()}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitCommit) Read(ctx context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	tflog.Debug(ctx, "Read git_add")
	// NO-OP: All data is already in Terraform state
}

func (r *resourceGitCommit) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update git_add")
	// NO-OP: All attributes require replacement, thus delete/create will be called
}

func (r *resourceGitCommit) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete git_add")
	// NO-OP: Terraform removes the state automatically for us
}
