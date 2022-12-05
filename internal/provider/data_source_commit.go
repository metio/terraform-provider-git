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

type CommitDataSource struct{}

var (
	_ datasource.DataSource           = (*CommitDataSource)(nil)
	_ datasource.DataSourceWithSchema = (*CommitDataSource)(nil)
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
	return &CommitDataSource{}
}

func (d *CommitDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_commit"
}

func (d *CommitDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches information about a single commit.",
		MarkdownDescription: "Fetches information about a single commit.",
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
				Description:         "The same value as the 'revision' attribute.",
				MarkdownDescription: "The same value as the `revision` attribute.",
				Computed:            true,
			},
			"revision": schema.StringAttribute{
				Description:         "The revision of the commit to fetch. Note that 'go-git' does not support every revision type at the moment. See https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision for details.",
				MarkdownDescription: "The [revision](https://www.git-scm.com/docs/gitrevisions) of the commit to fetch. Note that `go-git` does not [support](https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision) every revision type at the moment.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"sha1": schema.StringAttribute{
				Description:         "The SHA1 hash of the resolved revision.",
				MarkdownDescription: "The SHA1 hash of the resolved revision.",
				Computed:            true,
			},
			"author": schema.SingleNestedAttribute{
				Description:         "The original author of the commit.",
				MarkdownDescription: "The original author of the commit.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description:         "The name of the author.",
						MarkdownDescription: "The name of the author.",
						Computed:            true,
					},
					"email": schema.StringAttribute{
						Description:         "The email address of the author.",
						MarkdownDescription: "The email address of the author.",
						Computed:            true,
					},
					"timestamp": schema.StringAttribute{
						Description:         "The timestamp of the signature.",
						MarkdownDescription: "The timestamp of the signature.",
						Computed:            true,
					},
				},
			},
			"committer": schema.SingleNestedAttribute{
				Description:         "The person performing the commit.",
				MarkdownDescription: "The person performing the commit.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description:         "The name of the committer.",
						MarkdownDescription: "The name of the committer.",
						Computed:            true,
					},
					"email": schema.StringAttribute{
						Description:         "The email address of the committer.",
						MarkdownDescription: "The email address of the committer.",
						Computed:            true,
					},
					"timestamp": schema.StringAttribute{
						Description:         "The timestamp of the signature.",
						MarkdownDescription: "The timestamp of the signature.",
						Computed:            true,
					},
				},
			},
			"message": schema.StringAttribute{
				Description:         "The message of the commit.",
				MarkdownDescription: "The message of the commit.",
				Computed:            true,
			},
			"signature": schema.StringAttribute{
				Description:         "The signature of the commit.",
				MarkdownDescription: "The signature of the commit.",
				Computed:            true,
			},
			"tree_sha1": schema.StringAttribute{
				Description:         "The SHA1 checksum of the root tree of the commit.",
				MarkdownDescription: "The SHA1 checksum of the root tree of the commit.",
				Computed:            true,
			},
			"files": schema.ListAttribute{
				Description:         "The files updated by the commit.",
				MarkdownDescription: "The files updated by the commit.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *CommitDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_commit")

	var inputs commitDataSourceModel

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()
	revision := inputs.Revision.ValueString()

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
	state.SHA1 = types.StringValue(commitObject.Hash.String())
	state.Message = types.StringValue(commitObject.Message)
	state.Signature = types.StringValue(commitObject.PGPSignature)
	state.TreeSHA1 = types.StringValue(commitObject.TreeHash.String())
	state.Author = signatureToObject(&commitObject.Author)
	state.Committer = signatureToObject(&commitObject.Committer)
	state.Files, _ = types.ListValueFrom(ctx, types.StringType, extractModifiedFiles(commitObject))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
