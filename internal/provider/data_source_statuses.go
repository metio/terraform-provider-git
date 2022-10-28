/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type StatusesDataSource struct{}

var (
	_ datasource.DataSource = (*StatusesDataSource)(nil)
)

type statusesDataSourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	IsClean   types.Bool   `tfsdk:"is_clean"`
	Files     types.Map    `tfsdk:"files"`
}

func NewStatusesDataSource() datasource.DataSource {
	return &StatusesDataSource{}
}

func (d *StatusesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_statuses"
}

func (d *StatusesDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Fetches the status of all files in a Git repository.",
		MarkdownDescription: "Fetches the status of all files in a Git repository.",
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description:         "The path to the local Git repository.",
				MarkdownDescription: "The path to the local Git repository.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": {
				Description:         "The same value as the 'directory' attribute.",
				MarkdownDescription: "The same value as the `directory` attribute.",
				Type:                types.StringType,
				Computed:            true,
			},
			"is_clean": {
				Description:         "Whether the Git worktree is clean - all files must be in unmodified status for this to be true.",
				MarkdownDescription: "Whether the Git worktree is clean - all files must be in unmodified status for this to be true.",
				Type:                types.BoolType,
				Computed:            true,
			},
			"files": {
				Description:         "All modified files.",
				MarkdownDescription: "All modified files.",
				Computed:            true,
				Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
					"staging": {
						Description:         "The status of the file in the staging area.",
						MarkdownDescription: "The status of the file in the staging area.",
						Type:                types.StringType,
						Computed:            true,
					},
					"worktree": {
						Description:         "The status of the file in the worktree",
						MarkdownDescription: "The status of the file in the worktree",
						Type:                types.StringType,
						Computed:            true,
					},
				}),
			},
		},
	}, nil
}

func (d *StatusesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_statuses")

	var inputs statusesDataSourceModel
	var state statusesDataSourceModel

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	state.Directory = inputs.Directory
	state.Id = inputs.Directory

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
		state.IsClean = types.BoolValue(true)
		state.Files = types.MapValueMust(
			types.ObjectType{
				AttrTypes: statusType,
			},
			map[string]attr.Value{},
		)
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
		state.IsClean = types.BoolValue(status.IsClean())

		allFiles := make(map[string]attr.Value)
		for key, val := range status {
			allFiles[key] = types.ObjectValueMust(
				statusType,
				map[string]attr.Value{
					"staging":  types.String{Value: string(val.Staging)},
					"worktree": types.String{Value: string(val.Worktree)},
				},
			)
		}
		state.Files = types.MapValueMust(
			types.ObjectType{
				AttrTypes: statusType,
			},
			allFiles,
		)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
