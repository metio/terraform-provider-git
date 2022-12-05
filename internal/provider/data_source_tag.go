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

type TagDataSource struct{}

var (
	_ datasource.DataSource           = (*TagDataSource)(nil)
	_ datasource.DataSourceWithSchema = (*TagDataSource)(nil)
)

type tagDataSourceModel struct {
	Directory   types.String `tfsdk:"directory"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Lightweight types.Bool   `tfsdk:"lightweight"`
	Annotated   types.Bool   `tfsdk:"annotated"`
	SHA1        types.String `tfsdk:"sha1"`
	Message     types.String `tfsdk:"message"`
}

func NewTagDataSource() datasource.DataSource {
	return &TagDataSource{}
}

func (d *TagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (d *TagDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Reads information about a specific tag of a Git repository.",
		MarkdownDescription: "Reads information about a specific tag of a Git repository.",
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
				Description:         "The name of the tag to gather information about.",
				MarkdownDescription: "The name of the tag to gather information about.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"annotated": schema.BoolAttribute{
				Description:         "Whether the given tag is an annotated tag.",
				MarkdownDescription: "Whether the given tag is an annotated tag.",
				Computed:            true,
			},
			"lightweight": schema.BoolAttribute{
				Description:         "Whether the given tag is a lightweight tag.",
				MarkdownDescription: "Whether the given tag is a lightweight tag.",
				Computed:            true,
			},
			"sha1": schema.StringAttribute{
				Description:         "The SHA1 checksum of the commit the given tag is pointing at.",
				MarkdownDescription: "The SHA1 checksum of the commit the given tag is pointing at.",
				Computed:            true,
			},
			"message": schema.StringAttribute{
				Description:         "The associated message of an annotated tag.",
				MarkdownDescription: "The associated message of an annotated tag.",
				Computed:            true,
			},
		},
	}
}

func (d *TagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Git repository tag")

	var inputs tagDataSourceModel
	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()
	tagName := inputs.Name.ValueString()

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	tagReference := getTagReference(ctx, repository, tagName, &resp.Diagnostics)
	if tagReference == nil {
		return
	}

	tagObject, err := getTagObject(ctx, repository, tagReference.Hash(), &resp.Diagnostics)
	if err != nil {
		return
	}

	var state tagDataSourceModel
	state.Directory = inputs.Directory
	state.Id = inputs.Name
	state.Name = inputs.Name
	state.SHA1 = types.StringValue(tagReference.Hash().String())
	if tagObject == nil {
		state.Annotated = types.BoolValue(false)
		state.Lightweight = types.BoolValue(true)
		state.Message = types.StringNull()
	} else {
		state.Annotated = types.BoolValue(true)
		state.Lightweight = types.BoolValue(false)
		state.Message = types.StringValue(tagObject.Message)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
