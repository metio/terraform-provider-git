/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ConfigDataSource struct{}

var (
	_ datasource.DataSource = (*ConfigDataSource)(nil)
)

type configDataSourceModel struct {
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

func NewConfigDataSource() datasource.DataSource {
	return &ConfigDataSource{}
}

func (d *ConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (d *ConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Reads the configuration of a Git repository.",
		MarkdownDescription: "Reads the configuration of a Git repository.",
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
			"scope": schema.StringAttribute{
				Description:         "The configuration scope to read. Possible values are 'local', 'global', and 'system'. Defaults to 'global'.",
				MarkdownDescription: "The configuration scope to read. Possible values are `local`, `global`, and `system`. Defaults to `global`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"local",
						"global",
						"system",
					),
				},
			},
			"user_name": schema.StringAttribute{
				Description:         "The name of the author and the committer of a commit.",
				MarkdownDescription: "The name of the author and the committer of a commit.",
				Computed:            true,
			},
			"user_email": schema.StringAttribute{
				Description:         "The email address of the author and the committer of a commit.",
				MarkdownDescription: "The email address of the author and the committer of a commit.",
				Computed:            true,
			},
			"author_name": schema.StringAttribute{
				Description:         "The name of the author of a commit.",
				MarkdownDescription: "The name of the author of a commit.",
				Computed:            true,
			},
			"author_email": schema.StringAttribute{
				Description:         "The email address of the author of a commit.",
				MarkdownDescription: "The email address of the author of a commit.",
				Computed:            true,
			},
			"committer_name": schema.StringAttribute{
				Description:         "The name of the committer of a commit.",
				MarkdownDescription: "The name of the committer of a commit.",
				Computed:            true,
			},
			"committer_email": schema.StringAttribute{
				Description:         "The email address of the committer of a commit.",
				MarkdownDescription: "The email address of the committer of a commit.",
				Computed:            true,
			},
		},
	}
}

func (d *ConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_config")

	var inputs configDataSourceModel
	var state configDataSourceModel

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if inputs.Scope.IsNull() {
		inputs.Scope = types.StringValue("global")
	}

	directory := inputs.Directory.ValueString()

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	scope := inputs.Scope.ValueString()

	cfg := readConfig(ctx, repository, scope, &resp.Diagnostics)
	if cfg == nil {
		return
	}

	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.Scope = inputs.Scope
	state.UserName = types.StringValue(cfg.User.Name)
	state.UserEmail = types.StringValue(cfg.User.Email)
	state.AuthorName = types.StringValue(cfg.Author.Name)
	state.AuthorEmail = types.StringValue(cfg.Author.Email)
	state.CommitterName = types.StringValue(cfg.Committer.Name)
	state.CommitterEmail = types.StringValue(cfg.Committer.Email)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
