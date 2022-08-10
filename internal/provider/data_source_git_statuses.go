/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type dataSourceGitStatusesType struct{}

type dataSourceGitStatuses struct {
	p gitProvider
}

type dataSourceGitStatusesSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	IsClean   types.Bool   `tfsdk:"is_clean"`
	Files     types.Map    `tfsdk:"files"`
}

func (r *dataSourceGitStatusesType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Fetches the status of all files in a Git repository.",
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description: "The path to the local Git repository.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": {
				MarkdownDescription: "`DEPRECATED`: Only added in order to use the sdkv2 test framework. The path to the local Git repository.",
				Type:                types.StringType,
				Computed:            true,
			},
			"is_clean": {
				Description: "Whether the Git worktree is clean - all files must be in unmodified status for this to be true.",
				Type:        types.BoolType,
				Computed:    true,
			},
			"files": {
				Computed: true,
				Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
					"staging": {
						Description: "The status of the file in the staging area.",
						Type:        types.StringType,
						Computed:    true,
					},
					"worktree": {
						Description: "The status of the file in the worktree",
						Type:        types.StringType,
						Computed:    true,
					},
				}),
			},
		},
	}, nil
}

func (r *dataSourceGitStatusesType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return &dataSourceGitStatuses{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitStatuses) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	tflog.Debug(ctx, "Reading Git repository status")

	var inputs dataSourceGitStatusesSchema
	var state dataSourceGitStatusesSchema

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

	state.Directory = types.String{Value: directory}
	state.Id = types.String{Value: directory}

	statusType := map[string]attr.Type{
		"staging":  types.StringType,
		"worktree": types.StringType,
	}

	worktree, err := getWorktree(repository, &resp.Diagnostics)
	if err != nil {
		return
	} else if worktree == nil {
		tflog.Trace(ctx, "read worktree of bare repository", map[string]interface{}{
			"directory": directory,
		})
		state.IsClean = types.Bool{Value: true}
		state.Files = types.Map{
			ElemType: types.ObjectType{
				AttrTypes: statusType,
			},
			Elems: map[string]attr.Value{},
		}
	} else {
		tflog.Trace(ctx, "read worktree", map[string]interface{}{
			"directory": directory,
		})

		status, err := worktree.Status()
		if err != nil {
			resp.Diagnostics.AddError(
				"Cannot read status",
				"Could not read status because of: "+err.Error(),
			)
			return
		}
		tflog.Trace(ctx, "read status", map[string]interface{}{
			"directory": directory,
			"status":    status.String(),
		})
		state.IsClean = types.Bool{Value: status.IsClean()}

		allFiles := make(map[string]attr.Value)
		for key, val := range status {
			allFiles[key] = types.Object{
				AttrTypes: statusType,
				Attrs: map[string]attr.Value{
					"staging":  types.String{Value: string(val.Staging)},
					"worktree": types.String{Value: string(val.Worktree)},
				},
			}
		}
		state.Files = types.Map{
			ElemType: types.ObjectType{
				AttrTypes: statusType,
			},
			Elems: allFiles,
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
