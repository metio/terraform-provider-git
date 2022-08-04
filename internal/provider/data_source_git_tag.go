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
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func (r *dataSourceGitTagType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return &dataSourceGitTag{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitTag) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	tflog.Debug(ctx, "Reading Git repository tag")

	var config dataSourceGitTagSchema

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := config.Directory.Value
	tagName := config.Name.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	tagReference := getTagReference(ctx, repository, tagName, &resp.Diagnostics)
	if tagReference == nil {
		return
	}

	var state dataSourceGitTagSchema
	state.Directory = config.Directory
	state.Id = config.Name
	state.Name = config.Name
	state.SHA1 = types.String{Value: tagReference.Hash().String()}
	tag, err := repository.TagObject(tagReference.Hash())
	if err == plumbing.ErrObjectNotFound {
		state.Annotated = types.Bool{Value: false}
		state.Lightweight = types.Bool{Value: true}
		state.Message = types.String{Null: true}
	} else {
		state.Annotated = types.Bool{Value: true}
		state.Lightweight = types.Bool{Value: false}
		state.Message = types.String{Value: tag.Message}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
