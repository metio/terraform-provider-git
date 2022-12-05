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

type RemotesDataSource struct{}

var (
	_ datasource.DataSource           = (*RemotesDataSource)(nil)
	_ datasource.DataSourceWithSchema = (*RemotesDataSource)(nil)
)

type remotesDataSourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Remotes   types.Map    `tfsdk:"remotes"`
}

func NewRemotesDataSource() datasource.DataSource {
	return &RemotesDataSource{}
}

func (d *RemotesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remotes"
}

func (d *RemotesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Reads all configured remotes of a Git repository.",
		MarkdownDescription: "Reads all configured remotes of a Git repository.",
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
			"remotes": schema.MapNestedAttribute{
				Description:         "All configured remotes of the given Git repository.",
				MarkdownDescription: "All configured remotes of the given Git repository.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"urls": schema.ListAttribute{
							Description:         "The URLs for the remote.",
							MarkdownDescription: "The URLs for the remote.",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *RemotesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_remotes")

	var inputs remotesDataSourceModel
	var state remotesDataSourceModel

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

	remotes, err := repository.Remotes()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading remotes",
			"Could not read remotes of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "read remotes", map[string]interface{}{
		"directory": directory,
	})

	remoteType := map[string]attr.Type{
		"urls": types.ListType{ElemType: types.StringType},
	}

	allRemotes := make(map[string]attr.Value)
	for _, remote := range remotes {
		list, diags := types.ListValueFrom(ctx, types.StringType, remote.Config().URLs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		allRemotes[remote.Config().Name] = types.ObjectValueMust(
			remoteType,
			map[string]attr.Value{
				"urls": list,
			},
		)
	}

	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.Remotes = types.MapValueMust(
		types.ObjectType{
			AttrTypes: remoteType,
		},
		allRemotes,
	)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
