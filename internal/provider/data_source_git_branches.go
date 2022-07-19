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

type dataSourceGitBranchesType struct{}

type dataSourceGitBranches struct {
	p gitProvider
}

type dataSourceGitBranchesSchema struct {
	Directory types.String         `tfsdk:"directory"`
	Id        types.String         `tfsdk:"id"`
	Branches  map[string]GitBranch `tfsdk:"branches"`
}

func (r dataSourceGitBranchesType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Fetches all branches of a Git repository.",
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
			"branches": {
				Description: "All branches in a Git repository and their configuration.",
				Computed:    true,
				Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
					"sha1": {
						Description: "The SHA1 checksum of the HEAD of the branch.",
						Type:        types.StringType,
						Computed:    true,
					},
					"remote": {
						Description: "The name of remote this branch is tracking.",
						Type:        types.StringType,
						Computed:    true,
					},
					"rebase": {
						Description: "The rebase configuration of this branch.",
						Type:        types.StringType,
						Computed:    true,
					},
				}),
			},
		},
	}, nil
}

func (r dataSourceGitBranchesType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceGitBranches{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r dataSourceGitBranches) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var inputs dataSourceGitBranchesSchema
	var outputs dataSourceGitBranchesSchema

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

	branches, err := repository.Branches()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading branches",
			"Could not read branches of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "read branches", map[string]interface{}{
		"directory": directory,
	})

	allBranches := make(map[string]GitBranch)
	if err := branches.ForEach(func(reference *plumbing.Reference) error {
		branch, err := repository.Branch(reference.Name().Short())

		if branch != nil {
			allBranches[reference.Name().Short()] = GitBranch{
				SHA1:   types.String{Value: reference.Hash().String()},
				Remote: types.String{Value: branch.Remote},
				Rebase: types.String{Value: branch.Rebase},
			}
		}
		if err == git.ErrBranchNotFound {
			allBranches[reference.Name().Short()] = GitBranch{
				SHA1: types.String{Value: reference.Hash().String()},
			}
			return nil
		} else {
			return err
		}
	}); err != nil {
		resp.Diagnostics.AddError(
			"Error reading branches",
			"Could not read branches of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	outputs.Directory.Value = directory
	outputs.Id.Value = directory
	outputs.Branches = allBranches

	diags = resp.State.Set(ctx, &outputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
