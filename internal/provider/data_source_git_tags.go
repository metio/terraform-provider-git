/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
)

type dataSourceGitTagsType struct{}

type dataSourceGitTags struct {
	p gitProvider
}

type dataSourceGitTagsSchema struct {
	Directory   types.String `tfsdk:"directory"`
	Id          types.String `tfsdk:"id"`
	Lightweight types.Bool   `tfsdk:"lightweight"`
	Annotated   types.Bool   `tfsdk:"annotated"`
	Tags        types.Map    `tfsdk:"tags"`
}

func (r *dataSourceGitTagsType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Reads information about all tags of a Git repository.",
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
			"annotated": {
				MarkdownDescription: "Whether to request annotated tags. Defaults to `true`.",
				Type:                types.BoolType,
				Required:            false,
				Optional:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.Bool{Value: true}),
				},
			},
			"lightweight": {
				MarkdownDescription: "Whether to request lightweight tags. Defaults to `true`.",
				Type:                types.BoolType,
				Required:            false,
				Optional:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.Bool{Value: true}),
				},
			},
			"tags": {
				Computed: true,
				Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
					"annotated": {
						Description: "Whether the tag is an annotated tag or not.",
						Type:        types.BoolType,
						Computed:    true,
					},
					"lightweight": {
						Description: "Whether the tag is a lightweight tag or not.",
						Type:        types.BoolType,
						Computed:    true,
					},
					"sha1": {
						Description: "The SHA1 checksum of the commit the tag is pointing at.",
						Type:        types.StringType,
						Computed:    true,
					},
				}),
			},
		},
	}, nil
}

func (r *dataSourceGitTagsType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return &dataSourceGitTags{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitTags) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Git repository tags")

	var inputs dataSourceGitTagsSchema
	var state dataSourceGitTagsSchema

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
		inputs.Annotated = types.Bool{Value: true}
	}
	if inputs.Lightweight.IsNull() {
		inputs.Lightweight = types.Bool{Value: true}
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

		if inputs.Annotated.Value && tagObject != nil {
			allTags[ref.Name().Short()] = types.Object{
				AttrTypes: tagType,
				Attrs: map[string]attr.Value{
					"annotated":   types.Bool{Value: true},
					"lightweight": types.Bool{Value: false},
					"sha1":        types.String{Value: ref.Hash().String()},
				},
			}
		}
		if inputs.Lightweight.Value && tagObject == nil {
			allTags[ref.Name().Short()] = types.Object{
				AttrTypes: tagType,
				Attrs: map[string]attr.Value{
					"annotated":   types.Bool{Value: false},
					"lightweight": types.Bool{Value: true},
					"sha1":        types.String{Value: ref.Hash().String()},
				},
			}
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
	state.Tags = types.Map{
		ElemType: types.ObjectType{
			AttrTypes: tagType,
		},
		Elems: allTags,
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
