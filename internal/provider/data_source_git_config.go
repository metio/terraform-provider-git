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
	"github.com/metio/terraform-provider-git/internal/modifiers"
)

type dataSourceGitConfigType struct{}

type dataSourceGitConfig struct {
	p gitProvider
}

type dataSourceGitConfigSchema struct {
	Directory      types.String `tfsdk:"directory"`
	Id             types.String `tfsdk:"id"`
	Scope          types.String `tfsdk:"scope"`
	UserName       types.String `tfsdk:"user_name"`
	UserEmail      types.String `tfsdk:"user_email"`
	AuthorName     types.String `tfsdk:"author_name"`
	AuthorEmail    types.String `tfsdk:"author_email"`
	CommitterName  types.String `tfsdk:"committer_name"`
	CommitterEmail types.String `tfsdk:"committer_email"`
}

func (r *dataSourceGitConfigType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Reads the configuration of a Git repository.",
		MarkdownDescription: "Reads the configuration of a Git repository.",
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description:         "The path to the local Git repository.",
				MarkdownDescription: "The path to the local Git repository.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": {
				Description:         "The same value as the 'directory' attribute.",
				MarkdownDescription: "The same value as the `directory` attribute.",
				Type:                types.StringType,
				Computed:            true,
			},
			"scope": {
				Description:         "The configuration scope to read. Possible values are 'local', 'global', and 'system'. Defaults to 'global'.",
				MarkdownDescription: "The configuration scope to read. Possible values are `local`, `global`, and `system`. Defaults to `global`.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.String{Value: "global"}),
				},
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf(
						"local",
						"global",
						"system",
					),
				},
			},
			"user_name": {
				Description:         "The name of the author and the committer of a commit.",
				MarkdownDescription: "The name of the author and the committer of a commit.",
				Type:                types.StringType,
				Computed:            true,
			},
			"user_email": {
				Description:         "The email address of the author and the committer of a commit.",
				MarkdownDescription: "The email address of the author and the committer of a commit.",
				Type:                types.StringType,
				Computed:            true,
			},
			"author_name": {
				Description:         "The name of the author of a commit.",
				MarkdownDescription: "The name of the author of a commit.",
				Type:                types.StringType,
				Computed:            true,
			},
			"author_email": {
				Description:         "The email address of the author of a commit.",
				MarkdownDescription: "The email address of the author of a commit.",
				Type:                types.StringType,
				Computed:            true,
			},
			"committer_name": {
				Description:         "The name of the committer of a commit.",
				MarkdownDescription: "The name of the committer of a commit.",
				Type:                types.StringType,
				Computed:            true,
			},
			"committer_email": {
				Description:         "The email address of the committer of a commit.",
				MarkdownDescription: "The email address of the committer of a commit.",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (r *dataSourceGitConfigType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return &dataSourceGitConfig{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitConfig) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_config")

	var inputs dataSourceGitConfigSchema
	var state dataSourceGitConfigSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// NOTE: It seems default values for data sources are not working?
	if inputs.Scope.IsNull() {
		inputs.Scope = types.String{Value: "global"}
	}

	directory := inputs.Directory.Value
	scope := inputs.Scope.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	cfg, err := repository.ConfigScoped(mapConfigScope(scope))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading config",
			"Could not read git config because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "read config", map[string]interface{}{
		"directory": directory,
		"scope":     scope,
	})

	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.Scope = inputs.Scope
	state.UserName = types.String{Value: cfg.User.Name}
	state.UserEmail = types.String{Value: cfg.User.Email}
	state.AuthorName = types.String{Value: cfg.Author.Name}
	state.AuthorEmail = types.String{Value: cfg.Author.Email}
	state.CommitterName = types.String{Value: cfg.Committer.Name}
	state.CommitterEmail = types.String{Value: cfg.Committer.Email}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
