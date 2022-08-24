/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func New() provider.Provider {
	return &gitProvider{}
}

type gitProvider struct {
	configured bool
}

func (p *gitProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Provider for local [Git](https://git-scm.com/) operations. Requires Terraform 1.0 or later.",
		Attributes:          map[string]tfsdk.Attribute{},
	}, nil
}

func (p *gitProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
	p.configured = true
}

// GetResources - Defines gitProvider resources
func (p *gitProvider) GetResources(_ context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return map[string]provider.ResourceType{
		"git_init":   &resourceGitInitType{},
		"git_remote": &resourceGitRemoteType{},
		"git_tag":    &resourceGitTagType{},
	}, nil
}

// GetDataSources - Defines gitProvider data sources
func (p *gitProvider) GetDataSources(_ context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return map[string]provider.DataSourceType{
		"git_branch":     &dataSourceGitBranchType{},
		"git_branches":   &dataSourceGitBranchesType{},
		"git_commit":     &dataSourceGitCommitType{},
		"git_config":     &dataSourceGitConfigType{},
		"git_log":        &dataSourceGitLogType{},
		"git_remote":     &dataSourceGitRemoteType{},
		"git_remotes":    &dataSourceGitRemotesType{},
		"git_repository": &dataSourceGitRepositoryType{},
		"git_status":     &dataSourceGitStatusType{},
		"git_statuses":   &dataSourceGitStatusesType{},
		"git_tag":        &dataSourceGitTagType{},
		"git_tags":       &dataSourceGitTagsType{},
	}, nil
}
