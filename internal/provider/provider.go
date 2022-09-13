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

//// GetDataSources - Defines gitProvider data sources
//func (p *gitProvider) GetDataSources(_ context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
//	return map[string]provider.DataSourceType{
//		"git_branches":   &dataSourceGitBranchesType{},
//		"git_commit":     &dataSourceGitCommitType{},
//		"git_config":     &dataSourceGitConfigType{},
//		"git_log":        &dataSourceGitLogType{},
//		"git_remote":     &dataSourceGitRemoteType{},
//		"git_remotes":    &dataSourceGitRemotesType{},
//		"git_repository": &dataSourceGitRepositoryType{},
//		"git_status":     &dataSourceGitStatusType{},
//		"git_statuses":   &dataSourceGitStatusesType{},
//		"git_tag":        &dataSourceGitTagType{},
//		"git_tags":       &dataSourceGitTagsType{},
//	}, nil
//}
//
//// toProvider can be used to cast a generic provider.Provider reference to this specific provider.
//// This is ideally used in DataSourceType.NewDataSource and ResourceType.NewResource calls.
//func toProvider(in any) (*gitProvider, diag.Diagnostics) {
//	if in == nil {
//		return nil, nil
//	}
//
//	var diags diag.Diagnostics
//
//	p, ok := in.(*gitProvider)
//
//	if !ok {
//		diags.AddError(
//			"Unexpected Provider Instance Type",
//			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. "+
//				"This is always a bug in the provider code and should be reported to the provider developers.", in,
//			),
//		)
//		return nil, diags
//	}
//
//	return p, diags
//}
