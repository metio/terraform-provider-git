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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type commitDataSource struct{}

var (
	_ datasource.DataSource              = (*commitDataSource)(nil)
	_ datasource.DataSourceWithGetSchema = (*commitDataSource)(nil)
	_ datasource.DataSourceWithMetadata  = (*commitDataSource)(nil)
)

type commitDataSourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Revision  types.String `tfsdk:"revision"`
	SHA1      types.String `tfsdk:"sha1"`
	Author    types.Object `tfsdk:"author"`
	Committer types.Object `tfsdk:"committer"`
	Message   types.String `tfsdk:"message"`
	Signature types.String `tfsdk:"signature"`
	TreeSHA1  types.String `tfsdk:"tree_sha1"`
	Files     types.List   `tfsdk:"files"`
}

func NewCommitDataSource() datasource.DataSource {
	return &commitDataSource{}
}

func (d *commitDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_commit"
}

func (d *commitDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Fetches information about a single commit.",
		MarkdownDescription: "Fetches information about a single commit.",
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
				Description:         "The same value as the 'revision' attribute.",
				MarkdownDescription: "The same value as the `revision` attribute.",
				Type:                types.StringType,
				Computed:            true,
			},
			"revision": {
				Description:         "The revision of the commit to fetch. Note that 'go-git' does not support every revision type at the moment. See https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision for details.",
				MarkdownDescription: "The [revision](https://www.git-scm.com/docs/gitrevisions) of the commit to fetch. Note that `go-git` does not [support](https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision) every revision type at the moment.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"sha1": {
				Description:         "The SHA1 hash of the resolved revision.",
				MarkdownDescription: "The SHA1 hash of the resolved revision.",
				Type:                types.StringType,
				Computed:            true,
			},
			"author": {
				Description:         "The original author of the commit.",
				MarkdownDescription: "The original author of the commit.",
				Computed:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Description:         "The name of the author.",
						MarkdownDescription: "The name of the author.",
						Type:                types.StringType,
						Computed:            true,
					},
					"email": {
						Description:         "The email address of the author.",
						MarkdownDescription: "The email address of the author.",
						Type:                types.StringType,
						Computed:            true,
					},
					"timestamp": {
						Description:         "The timestamp of the signature.",
						MarkdownDescription: "The timestamp of the signature.",
						Type:                types.StringType,
						Computed:            true,
					},
				}),
			},
			"committer": {
				Description:         "The person performing the commit.",
				MarkdownDescription: "The person performing the commit.",
				Computed:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Description:         "The name of the committer.",
						MarkdownDescription: "The name of the committer.",
						Type:                types.StringType,
						Computed:            true,
					},
					"email": {
						Description:         "The email address of the committer.",
						MarkdownDescription: "The email address of the committer.",
						Type:                types.StringType,
						Computed:            true,
					},
					"timestamp": {
						Description:         "The timestamp of the signature.",
						MarkdownDescription: "The timestamp of the signature.",
						Type:                types.StringType,
						Computed:            true,
					},
				}),
			},
			"message": {
				Description:         "The message of the commit.",
				MarkdownDescription: "The message of the commit.",
				Type:                types.StringType,
				Computed:            true,
			},
			"signature": {
				Description:         "The signature of the commit.",
				MarkdownDescription: "The signature of the commit.",
				Type:                types.StringType,
				Computed:            true,
			},
			"tree_sha1": {
				Description:         "The SHA1 checksum of the root tree of the commit.",
				MarkdownDescription: "The SHA1 checksum of the root tree of the commit.",
				Type:                types.StringType,
				Computed:            true,
			},
			"files": {
				Description:         "The files updated by the commit.",
				MarkdownDescription: "The files updated by the commit.",
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Computed: true,
			},
		},
	}, nil
}

func (d *commitDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_commit")

	var inputs commitDataSourceModel

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

	hash := resolveRevision(ctx, repository, revision, &resp.Diagnostics)
	if hash == nil {
		return
	}

	commitObject := getCommit(ctx, repository, hash, &resp.Diagnostics)
	if commitObject == nil {
		return
	}

	var state commitDataSourceModel
	state.Directory = inputs.Directory
	state.Id = inputs.Revision
	state.Revision = inputs.Revision
	state.SHA1 = types.String{Value: commitObject.Hash.String()}
	state.Message = types.String{Value: commitObject.Message}
	state.Signature = types.String{Value: commitObject.PGPSignature}
	state.TreeSHA1 = types.String{Value: commitObject.TreeHash.String()}
	state.Author = signatureToObject(&commitObject.Author)
	state.Committer = signatureToObject(&commitObject.Committer)
	state.Files = stringsToList(extractModifiedFiles(commitObject))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
