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
				Description: "The current branch of the given Git repository. Note that repositories in detached state might not have a branch associated with them.",
				Type:        types.StringType,
				Computed:    true,
			},
			"sha1": {
				MarkdownDescription: "The current SHA1 of the `HEAD` of the given Git repository.",
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
	tflog.Debug(ctx, "Reading Git repository")

	var inputs dataSourceGitRepositorySchema
	var state dataSourceGitRepositorySchema

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
		"head":      head.String(),
	})

	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.SHA1 = types.String{Value: head.Hash().String()}
	if head.Name().IsBranch() {
		state.Branch = types.String{Value: head.Name().Short()}
	} else {
		state.Branch = types.String{Null: true}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
