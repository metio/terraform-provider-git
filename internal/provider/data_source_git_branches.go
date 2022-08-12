/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type dataSourceGitBranchesType struct{}

type dataSourceGitBranches struct {
	p gitProvider
}

type dataSourceGitBranchesSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Branches  types.Map    `tfsdk:"branches"`
}

func (r *dataSourceGitBranchesType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Fetches all branches of a Git repository.",
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

func (r *dataSourceGitBranchesType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return &dataSourceGitBranches{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitBranches) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Git repository branches")

	var inputs dataSourceGitBranchesSchema
	var state dataSourceGitBranchesSchema

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

	branchType := map[string]attr.Type{
		"sha1":   types.StringType,
		"remote": types.StringType,
		"rebase": types.StringType,
	}

	allBranches := make(map[string]attr.Value)
	if err := branches.ForEach(func(reference *plumbing.Reference) error {
		branch, err := repository.Branch(reference.Name().Short())

		if branch != nil {
			allBranches[reference.Name().Short()] = types.Object{
				AttrTypes: branchType,
				Attrs: map[string]attr.Value{
					"sha1":   types.String{Value: reference.Hash().String()},
					"remote": types.String{Value: branch.Remote},
					"rebase": types.String{Value: branch.Rebase},
				},
			}
		}
		if err == git.ErrBranchNotFound {
			allBranches[reference.Name().Short()] = types.Object{
				AttrTypes: branchType,
				Attrs: map[string]attr.Value{
					"sha1":   types.String{Value: reference.Hash().String()},
					"remote": types.String{Null: true},
					"rebase": types.String{Null: true},
				},
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

	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.Branches = types.Map{
		ElemType: types.ObjectType{
			AttrTypes: branchType,
		},
		Elems: allBranches,
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
