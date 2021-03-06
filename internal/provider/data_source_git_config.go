/*
 * This file is part of terraform-gitProvider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-gitProvider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5/config"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
)

type dataSourceGitConfigType struct{}

func (r dataSourceGitConfigType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Reads the configuration of a Git repository.",
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
			"scope": {
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
				Description: "The name of the author and the committer of a commit.",
				Type:        types.StringType,
				Computed:    true,
			},
			"user_email": {
				Description: "The email address of the author and the committer of a commit.",
				Type:        types.StringType,
				Computed:    true,
			},
			"author_name": {
				Description: "The name of the author of a commit.",
				Type:        types.StringType,
				Computed:    true,
			},
			"author_email": {
				Description: "The email address of the author of a commit.",
				Type:        types.StringType,
				Computed:    true,
			},
			"committer_name": {
				Description: "The name of the committer of a commit.",
				Type:        types.StringType,
				Computed:    true,
			},
			"committer_email": {
				Description: "The email address of the committer of a commit.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func (r dataSourceGitConfigType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceGitConfig{
		p: *(p.(*gitProvider)),
	}, nil
}

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

func (r dataSourceGitConfig) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var inputs dataSourceGitConfigSchema
	var outputs dataSourceGitConfigSchema

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

	repository := openRepository(ctx, directory, resp)
	if repository == nil {
		return
	}

	cfg, err := repository.ConfigScoped(mapScope(scope))
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

	outputs.Directory.Value = directory
	outputs.Id.Value = directory
	outputs.Scope.Value = scope
	outputs.UserName.Value = cfg.User.Name
	outputs.UserEmail.Value = cfg.User.Email
	outputs.AuthorName.Value = cfg.Author.Name
	outputs.AuthorEmail.Value = cfg.Author.Email
	outputs.CommitterName.Value = cfg.Committer.Name
	outputs.CommitterEmail.Value = cfg.Committer.Email

	diags = resp.State.Set(ctx, &outputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func mapScope(userInput string) config.Scope {
	switch userInput {
	case "local":
		return config.LocalScope
	case "system":
		return config.SystemScope
	default:
		return config.GlobalScope
	}
}
