/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type dataSourceGitTagType struct{}

type dataSourceGitTag struct {
	p gitProvider
}

type dataSourceGitTagSchema struct {
	Directory   types.String `tfsdk:"directory"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Lightweight types.Bool   `tfsdk:"lightweight"`
	Annotated   types.Bool   `tfsdk:"annotated"`
	SHA1        types.String `tfsdk:"sha1"`
	Message     types.String `tfsdk:"message"`
}

func (r *dataSourceGitTagType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Reads information about a specific tag of a Git repository.",
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
				MarkdownDescription: "`DEPRECATED`: Only added in order to use the sdkv2 test framework. The name of the tag to gather information about.",
				Type:                types.StringType,
				Computed:            true,
			},
			"name": {
				Description: "The name of the tag to gather information about.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"annotated": {
				Description: "Whether the given tag is an annotated tag.",
				Type:        types.BoolType,
				Computed:    true,
			},
			"lightweight": {
				Description: "Whether the given tag is a lightweight tag.",
				Type:        types.BoolType,
				Computed:    true,
			},
			"sha1": {
				Description: "The SHA1 checksum of the commit the given tag is pointing at.",
				Type:        types.StringType,
				Computed:    true,
			},
			"message": {
				Description: "The associated message of an annotated tag.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func (r *dataSourceGitTagType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return &dataSourceGitTag{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitTag) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Git repository tag")

	var inputs dataSourceGitTagSchema
	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	tagName := inputs.Name.Value

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

	var state dataSourceGitTagSchema
	state.Directory = inputs.Directory
	state.Id = inputs.Name
	state.Name = inputs.Name
	state.SHA1 = types.String{Value: tagReference.Hash().String()}
	if tagObject == nil {
		state.Annotated = types.Bool{Value: false}
		state.Lightweight = types.Bool{Value: true}
		state.Message = types.String{Null: true}
	} else {
		state.Annotated = types.Bool{Value: true}
		state.Lightweight = types.Bool{Value: false}
		state.Message = types.String{Value: tagObject.Message}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
