/*
 * SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
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
				MarkdownDescription: "`DEPRECATED`: Only added in order to use the sdkv2 test framework. The SHA1 checksum of the commit.",
				Type:                types.StringType,
				Computed:            true,
			},
			"sha1": {
				Description: "The SHA1 checksum of the commit.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(4),
				},
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
	tflog.Debug(ctx, "Reading Git commit information")

	var inputs dataSourceGitCommitSchema
	var state dataSourceGitCommitSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	sha1 := inputs.SHA1.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	state.Directory = inputs.Directory
	state.Id = inputs.SHA1
	state.SHA1 = inputs.SHA1

	commit, err := repository.CommitObject(plumbing.NewHash(sha1))
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read commit",
			"Could not read commit ["+sha1+"] because of: "+err.Error(),
		)
		return
	}

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
