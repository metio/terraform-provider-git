/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type GitProvider struct{}

var (
	_ provider.Provider = (*GitProvider)(nil)
)

func New() provider.Provider {
	return &GitProvider{}
}

func (p *GitProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "git"
}

func (p *GitProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Provider for local Git operations. Requires Terraform 1.0 or later.",
		MarkdownDescription: "Provider for local [Git](https://git-scm.com/) operations. Requires Terraform 1.0 or later.",
	}
}

func (p *GitProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
	// NO-OP: provider requires no configuration
}

func (p *GitProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewBranchDataSource,
		NewBranchesDataSource,
		NewCommitDataSource,
		NewConfigDataSource,
		NewLogDataSource,
		NewRemoteDataSource,
		NewRemotesDataSource,
		NewRepositoryDataSource,
		NewStatusDataSource,
		NewStatusesDataSource,
		NewTagDataSource,
		NewTagsDataSource,
	}
}

func (p *GitProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAddResource,
		NewCommitResource,
		NewInitResource,
		NewPushResource,
		NewRemoteResource,
		NewTagResource,
	}
}
