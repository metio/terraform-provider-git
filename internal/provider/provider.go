/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type gitProvider struct{}

var (
	_ provider.Provider                = (*gitProvider)(nil)
	_ provider.ProviderWithMetadata    = (*gitProvider)(nil)
	_ provider.ProviderWithDataSources = (*gitProvider)(nil)
	_ provider.ProviderWithResources   = (*gitProvider)(nil)
)

func New() provider.Provider {
	return &gitProvider{}
}

func (p *gitProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "git"
}

func (p *gitProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Provider for local Git operations. Requires Terraform 1.0 or later.",
		MarkdownDescription: "Provider for local [Git](https://git-scm.com/) operations. Requires Terraform 1.0 or later.",
	}, nil
}

func (p *gitProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
	// NO-OP: provider requires no configuration
}

func (p *gitProvider) DataSources(_ context.Context) []func() datasource.DataSource {
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

func (p *gitProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAddResource,
		NewCommitResource,
		NewInitResource,
		NewRemoteResource,
		NewTagResource,
	}
}
