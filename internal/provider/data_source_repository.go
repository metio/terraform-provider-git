/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type RepositoryDataSource struct{}

var (
	_ datasource.DataSource = (*RepositoryDataSource)(nil)
)

type repositoryDataSourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Branch    types.String `tfsdk:"branch"`
	SHA1      types.String `tfsdk:"sha1"`
}

func NewRepositoryDataSource() datasource.DataSource {
	return &RepositoryDataSource{}
}

func (d *RepositoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (d *RepositoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Reads information about a specific Git repository.",
		MarkdownDescription: "Reads information about a specific Git repository.",
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
			"branch": schema.StringAttribute{
				Description:         "The name of the current branch of the given Git repository. Note that repositories in detached state might not have a branch associated with them.",
				MarkdownDescription: "The name of the current branch of the given Git repository. Note that repositories in detached state might not have a branch associated with them.",
				Computed:            true,
			},
			"sha1": schema.StringAttribute{
				Description:         "The SHA1 of the current 'HEAD' of the given Git repository.",
				MarkdownDescription: "The SHA1 of the current `HEAD` of the given Git repository.",
				Computed:            true,
			},
		},
	}
}

func (d *RepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_repository")

	var inputs repositoryDataSourceModel

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

	var state repositoryDataSourceModel
	state.Directory = inputs.Directory
	state.Id = inputs.Directory

	head, err := repository.Head()
	if err == plumbing.ErrReferenceNotFound {
		state.SHA1 = types.StringNull()
		state.Branch = types.StringNull()
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error reading HEAD reference",
			"Could not read HEAD of ["+directory+"] because of: "+err.Error(),
		)
		return
	} else {
		tflog.Trace(ctx, "read HEAD reference", map[string]interface{}{
			"directory": directory,
			"head":      head.String(),
		})
		state.SHA1 = types.StringValue(head.Hash().String())
		if head.Name().IsBranch() {
			state.Branch = types.StringValue(head.Name().Short())
		} else {
			state.Branch = types.StringNull()
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
