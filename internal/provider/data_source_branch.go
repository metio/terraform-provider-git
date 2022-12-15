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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type BranchDataSource struct{}

var (
	_ datasource.DataSource = (*BranchDataSource)(nil)
)

type branchDataSourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	SHA1      types.String `tfsdk:"sha1"`
	Remote    types.String `tfsdk:"remote"`
	Rebase    types.String `tfsdk:"rebase"`
}

func NewBranchDataSource() datasource.DataSource {
	return &BranchDataSource{}
}

func (d *BranchDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch"
}

func (d *BranchDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches information about a specific branch of a Git repository.",
		MarkdownDescription: "Fetches information about a specific branch of a Git repository.",
		Attributes: map[string]schema.Attribute{
			"directory": schema.StringAttribute{
				Description:         "The path to the local Git repository.",
				MarkdownDescription: "The path to the local Git repository.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description:         "The name of the Git branch.",
				MarkdownDescription: "The name of the Git branch.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Description:         "The same value as the 'name' attribute.",
				MarkdownDescription: "The same value as the `name` attribute.",
				Computed:            true,
			},
			"sha1": schema.StringAttribute{
				Description:         "The SHA1 checksum of the 'HEAD' commit in the specified Git branch.",
				MarkdownDescription: "The SHA1 checksum of the `HEAD` commit in the specified Git branch.",
				Computed:            true,
			},
			"remote": schema.StringAttribute{
				Description:         "The configured remote for the specified Git branch.",
				MarkdownDescription: "The configured remote for the specified Git branch.",
				Computed:            true,
			},
			"rebase": schema.StringAttribute{
				Description:         "The rebase configuration for the specified Git branch. Possible values are 'true', 'interactive', and 'false'.",
				MarkdownDescription: "The rebase configuration for the specified Git branch. Possible values are `true`, `interactive`, and `false`.",
				Computed:            true,
			},
		},
	}
}

func (d *BranchDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_branch")

	var inputs branchDataSourceModel
	var state branchDataSourceModel

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()
	name := inputs.Name.ValueString()

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	branches, err := repository.Branches()
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read branches",
			"Could not read branches of ["+directory+"] because of: "+err.Error(),
		)
		return
	}
	state.SHA1 = types.StringUnknown()
	if err := branches.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().Short() == name {
			state.SHA1 = types.StringValue(ref.Hash().String())

			branch, err := repository.Branch(name)
			if branch != nil {
				state.Remote = types.StringValue(branch.Remote)
				state.Rebase = types.StringValue(branch.Rebase)
			} else if err == git.ErrBranchNotFound {
				state.Remote = types.StringNull()
				state.Rebase = types.StringNull()
			} else if err != nil {
				resp.Diagnostics.AddError(
					"Cannot read branch",
					"Could not read branch ["+name+"] of ["+directory+"] because of: "+err.Error(),
				)
				return err
			}

			tflog.Trace(ctx, "read branch", map[string]interface{}{
				"directory": directory,
				"branch":    name,
			})
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError(
			"Cannot read branches",
			"Could not read branches of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	if state.SHA1.IsNull() || state.SHA1.IsUnknown() {
		resp.Diagnostics.AddError(
			"Cannot read branch",
			"The branch ["+name+"] does not exist in ["+directory+"]",
		)
	}

	state.Directory = inputs.Directory
	state.Id = inputs.Name
	state.Name = inputs.Name

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
