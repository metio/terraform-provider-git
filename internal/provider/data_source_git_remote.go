/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
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

type dataSourceGitRemoteType struct{}

type dataSourceGitRemote struct {
	p gitProvider
}

type dataSourceGitRemoteSchema struct {
	Directory types.String   `tfsdk:"directory"`
	Id        types.String   `tfsdk:"id"`
	Remote    types.String   `tfsdk:"remote"`
	URLs      []types.String `tfsdk:"urls"`
}

func (r *dataSourceGitRemoteType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Reads information about a specific remote of a Git repository.",
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
			"remote": {
				Description: "The remote to gather information about.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"urls": {
				Description: "The configured URLs of the given remote.",
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Computed: true,
			},
		},
	}, nil
}

func (r *dataSourceGitRemoteType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return &dataSourceGitRemote{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitRemote) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	tflog.Debug(ctx, "Reading Git repository remote")

	var inputs dataSourceGitRemoteSchema
	var outputs dataSourceGitRemoteSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	requestedRemote := inputs.Remote.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	tflog.Trace(ctx, "opened repository", map[string]interface{}{
		"directory": directory,
	})

	remote, err := repository.Remote(requestedRemote)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read remote",
			"Could not read remote ["+requestedRemote+"] of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "read remote", map[string]interface{}{
		"directory": directory,
		"remote":    requestedRemote,
	})

	outputs.Directory = types.String{Value: directory}
	outputs.Id = types.String{Value: directory}
	outputs.Remote = types.String{Value: remote.Config().Name}
	outputs.URLs = extractGitRemoteUrls(remote)

	diags = resp.State.Set(ctx, &outputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
