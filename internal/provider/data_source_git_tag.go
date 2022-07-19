/*
 * This file is part of terraform-gitProvider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-gitProvider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
	Tag         types.String `tfsdk:"tag"`
	Lightweight types.Bool   `tfsdk:"lightweight"`
	Annotated   types.Bool   `tfsdk:"annotated"`
	SHA1        types.String `tfsdk:"sha1"`
}

func (r dataSourceGitTagType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Reads information about a specific tag of a Git repository.",
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description: "The path to the local Git repository.",
				Type:        types.StringType,
				Required:    true,
			},
			"id": {
				MarkdownDescription: "`DEPRECATED`: Only added in order to use the sdkv2 test framework. The path to the local Git repository.",
				Type:                types.StringType,
				Computed:            true,
			},
			"tag": {
				Description: "The tag to gather information about.",
				Type:        types.StringType,
				Required:    true,
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
		},
	}, nil
}

func (r dataSourceGitTagType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceGitTag{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r dataSourceGitTag) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var inputs dataSourceGitTagSchema
	var outputs dataSourceGitTagSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	repository, err := git.PlainOpen(directory)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error opening repository",
			"Could not open git repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "opened repository", map[string]interface{}{
		"directory": directory,
	})

	tag := inputs.Tag.Value
	reference, err := repository.Tag(tag)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tag",
			"Could not read tag ["+tag+"] of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "read tag", map[string]interface{}{
		"directory": directory,
		"tag":       tag,
	})

	outputs.Directory.Value = directory
	outputs.Id.Value = directory
	outputs.Tag.Value = reference.Name().Short()
	outputs.SHA1.Value = reference.Hash().String()
	_, err = repository.TagObject(reference.Hash())
	if err == plumbing.ErrObjectNotFound {
		outputs.Annotated.Value = false
		outputs.Lightweight.Value = true
	} else {
		outputs.Annotated.Value = true
		outputs.Lightweight.Value = false
	}

	diags = resp.State.Set(ctx, &outputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
