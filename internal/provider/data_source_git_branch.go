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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type dataSourceGitBranchType struct{}

type dataSourceGitBranch struct {
	p gitProvider
}

type dataSourceGitBranchSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Branch    types.String `tfsdk:"branch"`
	SHA1      types.String `tfsdk:"sha1"`
	Remote    types.String `tfsdk:"remote"`
	Rebase    types.String `tfsdk:"rebase"`
}

func (r *dataSourceGitBranchType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Fetches information about a specific branch of a Git repository.",
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description: "The path to the local Git repository.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"branch": {
				Description: "The name of the Git branch.",
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
			"sha1": {
				Description: "The SHA1 checksum of the HEAD commit in the specified Git branch.",
				Type:        types.StringType,
				Computed:    true,
			},
			"remote": {
				Description: "The configured remote for the specified Git branch.",
				Type:        types.StringType,
				Computed:    true,
			},
			"rebase": {
				Description: "The rebase configuration for the specified Git branch.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func (r *dataSourceGitBranchType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return &dataSourceGitBranch{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitBranch) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Git repository branch")

	var inputs dataSourceGitBranchSchema
	var state dataSourceGitBranchSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	requestedBranch := inputs.Branch.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	branch, err := repository.Branch(requestedBranch)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read branch",
			"Could not read branch ["+requestedBranch+"] of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "read branch", map[string]interface{}{
		"directory": directory,
		"branch":    branch.Name,
	})

	branches, err := repository.Branches()
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read branches",
			"Could not read branches of ["+directory+"] because of: "+err.Error(),
		)
		return
	}
	if err := branches.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().Short() == branch.Name {
			state.SHA1 = types.String{Value: ref.Hash().String()}
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError(
			"Cannot read branches",
			"Could not read branches of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.Branch = inputs.Branch
	state.Remote = types.String{Value: branch.Remote}
	state.Rebase = types.String{Value: branch.Rebase}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
