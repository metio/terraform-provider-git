/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type RemoteDataSource struct{}

var (
	_ datasource.DataSource = (*RemoteDataSource)(nil)
)

type remoteDataSourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	URLs      types.List   `tfsdk:"urls"`
}

func NewRemoteDataSource() datasource.DataSource {
	return &RemoteDataSource{}
}

func (d *RemoteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote"
}

func (d *RemoteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Reads information about a specific remote of a Git repository.",
		MarkdownDescription: "Reads information about a specific remote of a Git repository.",
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
				Description:         "The same value as the 'name' attribute.",
				MarkdownDescription: "The same value as the `name` attribute.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "The name of the remote to gather information about.",
				MarkdownDescription: "The name of the remote to gather information about.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"urls": schema.ListAttribute{
				Description:         "The configured URLs of the given remote.",
				MarkdownDescription: "The configured URLs of the given remote.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *RemoteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_remote")

	var inputs remoteDataSourceModel
	var state remoteDataSourceModel

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()
	remoteName := inputs.Name.ValueString()

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	remote := getRemote(ctx, repository, remoteName, &resp.Diagnostics)
	if remote == nil {
		return
	}

	state.Directory = inputs.Directory
	state.Id = inputs.Name
	state.Name = inputs.Name
	state.URLs, _ = types.ListValueFrom(ctx, types.StringType, remote.Config().URLs)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
