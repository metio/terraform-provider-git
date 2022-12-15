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
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

func (d *StatusesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the status of all files in a Git repository.",
		MarkdownDescription: "Fetches the status of all files in a Git repository.",
		Attributes: map[string]schema.Attribute{
			"directory": schema.StringAttribute{
				Description:         "The path to the local Git repository.",
				MarkdownDescription: "The path to the local Git repository.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Description:         "The same value as the 'directory' attribute.",
				MarkdownDescription: "The same value as the `directory` attribute.",
				Computed:            true,
			},
			"is_clean": schema.BoolAttribute{
				Description:         "Whether the Git worktree is clean - all files must be in unmodified status for this to be true.",
				MarkdownDescription: "Whether the Git worktree is clean - all files must be in unmodified status for this to be true.",
				Computed:            true,
			},
			"files": schema.MapNestedAttribute{
				Description:         "All modified files.",
				MarkdownDescription: "All modified files.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"staging": schema.StringAttribute{
							Description:         "The status of the file in the staging area.",
							MarkdownDescription: "The status of the file in the staging area.",
							Computed:            true,
						},
						"worktree": schema.StringAttribute{
							Description:         "The status of the file in the worktree",
							MarkdownDescription: "The status of the file in the worktree",
							Computed:            true,
						},
					},
				},
			},
		},
	}
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
					"staging":  types.StringValue(string(val.Staging)),
					"worktree": types.StringValue(string(val.Worktree)),
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
