/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type StatusDataSource struct{}

var (
	_ datasource.DataSource = (*StatusDataSource)(nil)
)

type statusDataSourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	File      types.String `tfsdk:"file"`
	Staging   types.String `tfsdk:"staging"`
	Worktree  types.String `tfsdk:"worktree"`
}

func NewStatusDataSource() datasource.DataSource {
	return &StatusDataSource{}
}

func (d *StatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status"
}

func (d *StatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the status of a single files in a Git repository.",
		MarkdownDescription: "Fetches the status of a single files in a Git repository.",
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
				Description:         "The same value as the 'file' attribute.",
				MarkdownDescription: "The same value as the `file` attribute.",
				Computed:            true,
			},
			"file": schema.StringAttribute{
				Description:         "The file to get status information about.",
				MarkdownDescription: "The file to get status information about.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
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
	}
}

func (d *StatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_status")

	var inputs statusDataSourceModel
	var state statusDataSourceModel

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()
	fileName := inputs.File.ValueString()

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	state.Directory = inputs.Directory
	state.Id = inputs.File
	state.File = inputs.File
	state.Staging = types.StringNull()
	state.Worktree = types.StringNull()

	worktree, err := repository.Worktree()
	if errors.Is(err, git.ErrIsBareRepository) {
		tflog.Trace(ctx, "read worktree of bare repository", map[string]interface{}{
			"directory": directory,
		})
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read worktree",
			"Could not read worktree because of: "+err.Error(),
		)
		return
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

		fileStatus := status.File(fileName)
		tflog.Trace(ctx, "read file status", map[string]interface{}{
			"directory": directory,
			"file":      fileName,
			"staging":   fileStatus.Staging,
			"worktree":  fileStatus.Worktree,
		})
		state.Staging = types.StringValue(string(fileStatus.Staging))
		state.Worktree = types.StringValue(string(fileStatus.Worktree))
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
