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

type dataSourceGitRepositoryType struct{}

type dataSourceGitRepository struct {
	p gitProvider
}

type dataSourceGitRepositorySchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Branch    types.String `tfsdk:"branch"`
}

func (r *dataSourceGitRepositoryType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Reads information about a specific Git repository.",
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
			"branch": {
				Description: "The current branch of the given Git repository.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func (r *dataSourceGitRepositoryType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return &dataSourceGitRepository{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitRepository) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	tflog.Debug(ctx, "Reading Git repository")

	var inputs dataSourceGitRepositorySchema
	var outputs dataSourceGitRepositorySchema

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

	head, err := repository.Head()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading HEAD reference",
			"Could not read HEAD of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "read HEAD reference", map[string]interface{}{
		"directory": directory,
		"head":      head.Name().String(),
	})

	outputs.Directory.Value = directory
	outputs.Id.Value = directory
	outputs.Branch.Value = head.Name().Short()

	diags = resp.State.Set(ctx, &outputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
