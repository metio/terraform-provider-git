/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type dataSourceGitCommitType struct{}

type dataSourceGitCommit struct {
	p gitProvider
}

type dataSourceGitCommitSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Revision  types.String `tfsdk:"revision"`
	SHA1      types.String `tfsdk:"sha1"`
	Author    types.Object `tfsdk:"author"`
	Committer types.Object `tfsdk:"committer"`
	Message   types.String `tfsdk:"message"`
	Signature types.String `tfsdk:"signature"`
	TreeSHA1  types.String `tfsdk:"tree_sha1"`
}

func (r *dataSourceGitCommitType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Fetches information about a single commit.",
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
				MarkdownDescription: "The same value as the `revision` attribute.",
				Type:                types.StringType,
				Computed:            true,
			},
			"revision": {
				MarkdownDescription: "The [revision](https://www.git-scm.com/docs/gitrevisions) of the commit to fetch. Note that `go-git` does not [support](https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision) every revision type at the moment.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"sha1": {
				MarkdownDescription: "The SHA1 hash of the resolved revision.",
				Type:                types.StringType,
				Computed:            true,
			},
			"author": {
				Description: "The original author of the commit.",
				Computed:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Description: "The name of the author.",
						Type:        types.StringType,
						Computed:    true,
					},
					"email": {
						Description: "The email address of the author.",
						Type:        types.StringType,
						Computed:    true,
					},
					"timestamp": {
						Description: "The timestamp of the signature.",
						Type:        types.StringType,
						Computed:    true,
					},
				}),
			},
			"committer": {
				Description: "The person performing the commit.",
				Computed:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Description: "The name of the committer.",
						Type:        types.StringType,
						Computed:    true,
					},
					"email": {
						Description: "The email address of the committer.",
						Type:        types.StringType,
						Computed:    true,
					},
					"timestamp": {
						Description: "The timestamp of the signature.",
						Type:        types.StringType,
						Computed:    true,
					},
				}),
			},
			"message": {
				Description: "The message of the commit.",
				Type:        types.StringType,
				Computed:    true,
			},
			"signature": {
				Description: "The signature of the commit.",
				Type:        types.StringType,
				Computed:    true,
			},
			"tree_sha1": {
				Description: "The SHA1 checksum of the root tree of the commit.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func (r *dataSourceGitCommitType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return &dataSourceGitCommit{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitCommit) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_commit")

	var inputs dataSourceGitCommitSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	revision := inputs.Revision.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	hash, err := repository.ResolveRevision(plumbing.Revision(revision))
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot resolve revision",
			"Could not resolve revision ["+revision+"] because of: "+err.Error(),
		)
		return
	}

	commit, err := repository.CommitObject(*hash)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read commit",
			"Could not read commit ["+hash.String()+"] because of: "+err.Error(),
		)
		return
	}

	var state dataSourceGitCommitSchema
	state.Directory = inputs.Directory
	state.Id = inputs.Revision
	state.Revision = inputs.Revision
	state.SHA1 = types.String{Value: commit.Hash.String()}
	state.Message = types.String{Value: commit.Message}
	state.Signature = types.String{Value: commit.PGPSignature}
	state.TreeSHA1 = types.String{Value: commit.TreeHash.String()}
	state.Author = signatureToObject(&commit.Author)
	state.Committer = signatureToObject(&commit.Committer)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
