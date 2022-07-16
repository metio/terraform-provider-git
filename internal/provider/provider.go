/*
 * This file is part of terraform-gitProvider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-gitProvider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func New() tfsdk.Provider {
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

func (p *gitProvider) Configure(_ context.Context, _ tfsdk.ConfigureProviderRequest, _ *tfsdk.ConfigureProviderResponse) {
	p.configured = true
}

// GetResources - Defines gitProvider resources
func (p *gitProvider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		//"git_init": resourceGitInit{},
	}, nil
}

// GetDataSources - Defines gitProvider data sources
func (p *gitProvider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"git_branch":     dataSourceGitBranchType{},
		"git_branches":   dataSourceGitBranchesType{},
		"git_config":     dataSourceGitConfigType{},
		"git_repository": dataSourceGitRepositoryType{},
		"git_remote":     dataSourceGitRemoteType{},
		"git_remotes":    dataSourceGitRemotesType{},
		"git_tag":        dataSourceGitTagType{},
		"git_tags":       dataSourceGitTagsType{},
	}, nil
}
