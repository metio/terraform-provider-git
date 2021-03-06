/*
 * This file is part of terraform-gitProvider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-gitProvider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type dataSourceGitRemotesType struct{}

type dataSourceGitRemotes struct {
	p gitProvider
}

type dataSourceGitRemotesSchema struct {
	Directory types.String         `tfsdk:"directory"`
	Id        types.String         `tfsdk:"id"`
	Remotes   map[string]GitRemote `tfsdk:"remotes"`
}

func (r dataSourceGitRemotesType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Reads all configured remotes of a Git repository.",
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
			"remotes": {
				Description: "All configured remotes of the given Git repository.",
				Computed:    true,
				Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
					"urls": {
						Description: "The URLs for the remote.",
						Type: types.ListType{
							ElemType: types.StringType,
						},
						Computed: true,
					},
				}),
			},
		},
	}, nil
}

func (r dataSourceGitRemotesType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceGitRemotes{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r dataSourceGitRemotes) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var inputs dataSourceGitRemotesSchema
	var outputs dataSourceGitRemotesSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value

	repository := openRepository(ctx, directory, resp)
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

	allRemotes := make(map[string]GitRemote)
	for _, remote := range remotes {
		allRemotes[remote.Config().Name] = GitRemote{
			URLs: extractGitRemoteUrls(remote),
		}
	}

	outputs.Directory.Value = directory
	outputs.Id.Value = directory
	outputs.Remotes = allRemotes

	diags = resp.State.Set(ctx, &outputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
