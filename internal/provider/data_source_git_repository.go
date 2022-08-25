/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
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

type dataSourceGitRepositoryType struct{}

type dataSourceGitRepository struct {
	p gitProvider
}

type dataSourceGitRepositorySchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Branch    types.String `tfsdk:"branch"`
	SHA1      types.String `tfsdk:"sha1"`
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
				MarkdownDescription: "The same value as the `directory` attribute.",
				Type:                types.StringType,
				Computed:            true,
			},
			"branch": {
				Description: "The name of the current branch of the given Git repository. Note that repositories in detached state might not have a branch associated with them.",
				Type:        types.StringType,
				Computed:    true,
			},
			"sha1": {
				MarkdownDescription: "The SHA1 of the current `HEAD` of the given Git repository.",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (r *dataSourceGitRepositoryType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return &dataSourceGitRepository{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitRepository) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read git_repository")

	var inputs dataSourceGitRepositorySchema

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

	var state dataSourceGitRepositorySchema
	state.Directory = inputs.Directory
	state.Id = inputs.Directory

	head, err := repository.Head()
	if err == plumbing.ErrReferenceNotFound {
		state.SHA1 = types.String{Null: true}
		state.Branch = types.String{Null: true}
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error reading HEAD reference",
			"Could not read HEAD of ["+directory+"] because of: "+err.Error(),
		)
		return
	} else {
		tflog.Trace(ctx, "read HEAD reference", map[string]interface{}{
			"directory": directory,
			"head":      head.String(),
		})
		state.SHA1 = types.String{Value: head.Hash().String()}
		if head.Name().IsBranch() {
			state.Branch = types.String{Value: head.Name().Short()}
		} else {
			state.Branch = types.String{Null: true}
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
