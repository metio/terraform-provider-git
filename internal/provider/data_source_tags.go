/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TagsDataSource struct{}

var (
	_ datasource.DataSource           = (*TagsDataSource)(nil)
	_ datasource.DataSourceWithSchema = (*TagsDataSource)(nil)
)

type tagsDataSourceModel struct {
	Directory   types.String `tfsdk:"directory"`
	Id          types.String `tfsdk:"id"`
	Lightweight types.Bool   `tfsdk:"lightweight"`
	Annotated   types.Bool   `tfsdk:"annotated"`
	Tags        types.Map    `tfsdk:"tags"`
}

func NewTagsDataSource() datasource.DataSource {
	return &TagsDataSource{}
}

func (d *TagsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tags"
}

func (d *TagsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Reads information about all tags of a Git repository.",
		MarkdownDescription: "Reads information about all tags of a Git repository.",
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
			"annotated": schema.BoolAttribute{
				Description:         "Whether to request annotated tags. Defaults to 'true'.",
				MarkdownDescription: "Whether to request annotated tags. Defaults to `true`.",
				Required:            false,
				Optional:            true,
			},
			"lightweight": schema.BoolAttribute{
				Description:         "Whether to request lightweight tags. Defaults to 'true'.",
				MarkdownDescription: "Whether to request lightweight tags. Defaults to `true`.",
				Required:            false,
				Optional:            true,
			},
			"tags": schema.MapNestedAttribute{
				Description:         "All existing tags.",
				MarkdownDescription: "All existing tags.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"annotated": schema.BoolAttribute{
							Description:         "Whether the tag is an annotated tag or not.",
							MarkdownDescription: "Whether the tag is an annotated tag or not.",
							Computed:            true,
						},
						"lightweight": schema.BoolAttribute{
							Description:         "Whether the tag is a lightweight tag or not.",
							MarkdownDescription: "Whether the tag is a lightweight tag or not.",
							Computed:            true,
						},
						"sha1": schema.StringAttribute{
							Description:         "The SHA1 checksum of the commit the tag is pointing at.",
							MarkdownDescription: "The SHA1 checksum of the commit the tag is pointing at.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *TagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_tags")

	var inputs tagsDataSourceModel
	var state tagsDataSourceModel

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

	tags, err := repository.Tags()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tags",
			"Could not read tags of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "read tags", map[string]interface{}{
		"directory": directory,
	})

	// NOTE: It seems default values for data sources are not working?
	if inputs.Annotated.IsNull() {
		inputs.Annotated = types.BoolValue(true)
	}
	if inputs.Lightweight.IsNull() {
		inputs.Lightweight = types.BoolValue(true)
	}

	tagType := map[string]attr.Type{
		"annotated":   types.BoolType,
		"lightweight": types.BoolType,
		"sha1":        types.StringType,
	}

	allTags := make(map[string]attr.Value)
	if err := tags.ForEach(func(ref *plumbing.Reference) error {
		tagObject, err := getTagObject(ctx, repository, ref.Hash(), &resp.Diagnostics)
		if err != nil {
			return err
		}

		if inputs.Annotated.ValueBool() && tagObject != nil {
			allTags[ref.Name().Short()] = types.ObjectValueMust(
				tagType,
				map[string]attr.Value{
					"annotated":   types.BoolValue(true),
					"lightweight": types.BoolValue(false),
					"sha1":        types.StringValue(ref.Hash().String()),
				},
			)
		}
		if inputs.Lightweight.ValueBool() && tagObject == nil {
			allTags[ref.Name().Short()] = types.ObjectValueMust(
				tagType,
				map[string]attr.Value{
					"annotated":   types.BoolValue(false),
					"lightweight": types.BoolValue(true),
					"sha1":        types.StringValue(ref.Hash().String()),
				},
			)
		}

		return nil
	}); err != nil {
		resp.Diagnostics.AddError(
			"Error reading tags",
			"Could not read tags of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.Annotated = inputs.Annotated
	state.Lightweight = inputs.Lightweight
	state.Tags = types.MapValueMust(
		types.ObjectType{
			AttrTypes: tagType,
		},
		allTags,
	)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
