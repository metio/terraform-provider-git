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
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func (c *resourceGitInitType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Initializes a Git repository.",
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
			"bare": {
				Description: "Whether the created Git repository is bare or not. Defaults to `false`.",
				Type:        types.BoolType,
				Computed:    true,
				Optional:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.Bool{Value: false}),
					tfsdk.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (r *resourceGitInitType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return &resourceGitInit{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *resourceGitInit) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, "Creating Git repository")

	var inputs resourceGitInitSchema
	var state resourceGitInitSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	bare := inputs.Bare.Value

	_, err := git.PlainInit(directory, bare)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot create repository",
			"Could not create repository ["+directory+"] because of: "+err.Error(),
		)
	}

	tflog.Trace(ctx, "created repository", map[string]interface{}{
		"directory": directory,
		"bare":      bare,
	})

	state.Directory = types.String{Value: directory}
	state.Id = types.String{Value: directory}
	state.Bare = types.Bool{Value: bare}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitInit) Read(ctx context.Context, _ tfsdk.ReadResourceRequest, _ *tfsdk.ReadResourceResponse) {
	tflog.Debug(ctx, "Reading Git repository")
}

func (r *resourceGitInit) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "Updating Git repository")
	updatedUsingPlan(ctx, &req, resp, &resourceGitInitSchema{})
}

func (r *resourceGitInit) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "Removing Git repository")

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

func (r *resourceGitInit) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
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
	_, err := repository.Worktree()
	if err == git.ErrIsBareRepository {
		state.Bare = types.Bool{Value: true}
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read worktree",
			"Could not read worktree because of: "+err.Error(),
		)
		return
	} else {
		state.Bare = types.Bool{Value: false}
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
