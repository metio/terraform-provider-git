/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type BranchesDataSource struct{}

var (
	_ datasource.DataSource = (*BranchesDataSource)(nil)
)

type branchesDataSourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Branches  types.Map    `tfsdk:"branches"`
}

func NewBranchesDataSource() datasource.DataSource {
	return &BranchesDataSource{}
}

func (d *BranchesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branches"
}

func (d *BranchesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches all branches of a Git repository.",
		MarkdownDescription: "Fetches all branches of a Git repository.",
		Attributes: map[string]schema.Attribute{
			"directory": schema.StringAttribute{
				Description:         "The path to the local Git repository.",
				MarkdownDescription: "The path to the local Git repository.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Description:         "The same value as the 'directory' attribute.",
				MarkdownDescription: "The same value as the `directory` attribute.",
				Computed:            true,
			},
			"branches": schema.MapNestedAttribute{
				Description:         "All branches in a Git repository and their configuration.",
				MarkdownDescription: "All branches in a Git repository and their configuration.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sha1": schema.StringAttribute{
							Description:         "The SHA1 checksum of the 'HEAD' of the branch.",
							MarkdownDescription: "The SHA1 checksum of the `HEAD` of the branch.",
							Computed:            true,
						},
						"remote": schema.StringAttribute{
							Description:         "The name of remote this branch is tracking.",
							MarkdownDescription: "The name of remote this branch is tracking.",
							Computed:            true,
						},
						"rebase": schema.StringAttribute{
							Description:         "The rebase configuration for the specified Git branch. Possible values are 'true', 'interactive', and 'false'.",
							MarkdownDescription: "The rebase configuration for the specified Git branch. Possible values are `true`, `interactive`, and `false`.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *BranchesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_branches")

	var inputs branchesDataSourceModel
	var state branchesDataSourceModel

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()

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
			allBranches[reference.Name().Short()] = types.ObjectValueMust(
				branchType,
				map[string]attr.Value{
					"sha1":   types.StringValue(reference.Hash().String()),
					"remote": types.StringValue(branch.Remote),
					"rebase": types.StringValue(branch.Rebase),
				},
			)
		}
		if err == git.ErrBranchNotFound {
			allBranches[reference.Name().Short()] = types.ObjectValueMust(
				branchType,
				map[string]attr.Value{
					"sha1":   types.StringValue(reference.Hash().String()),
					"remote": types.StringNull(),
					"rebase": types.StringNull(),
				},
			)
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
	state.Branches = types.MapValueMust(
		types.ObjectType{
			AttrTypes: branchType,
		},
		allBranches,
	)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
