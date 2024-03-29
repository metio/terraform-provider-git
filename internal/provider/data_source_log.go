/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type LogDataSource struct{}

var (
	_ datasource.DataSource = (*LogDataSource)(nil)
)

type logDataSourceModel struct {
	Directory   types.String `tfsdk:"directory"`
	Id          types.String `tfsdk:"id"`
	From        types.String `tfsdk:"from"`
	Order       types.String `tfsdk:"order"`
	All         types.Bool   `tfsdk:"all"`
	Since       types.String `tfsdk:"since"`
	Until       types.String `tfsdk:"until"`
	MaxCount    types.Int64  `tfsdk:"max_count"`
	Skip        types.Int64  `tfsdk:"skip"`
	FilterPaths types.List   `tfsdk:"filter_paths"`
	Commits     types.List   `tfsdk:"commits"`
}

func NewLogDataSource() datasource.DataSource {
	return &LogDataSource{}
}

func (d *LogDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log"
}

func (d *LogDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the commit log of a Git repository similar to 'git log'.",
		MarkdownDescription: "Fetches the commit log of a Git repository similar to `git log`.",
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
			"from": schema.StringAttribute{
				Description:         "When set the log will only contain commits reachable from it. If this option is not set, 'HEAD' will be used as the default. Can be any revision that 'go-git' supports. See https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision for details.",
				MarkdownDescription: "When set the log will only contain commits reachable from it. If this option is not set, `HEAD` will be used as the default. Can be any [revision](https://www.git-scm.com/docs/gitrevisions) that `go-git` [supports](https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision).",
				Optional:            true,
				Computed:            true,
			},
			"order": schema.StringAttribute{
				Description:         "The traversal algorithm to use while listing commits. Defaults to 'time' which is similar to 'git log'. Other values are 'depth' and 'breadth' for depth- or breadth-first traversal.",
				MarkdownDescription: "The traversal algorithm to use while listing commits. Defaults to `time` which is similar to `git log`. Other values are `depth` and `breadth` for depth- or breadth-first traversal.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"time",
						"depth",
						"breadth",
					),
				},
			},
			"all": schema.BoolAttribute{
				Description:         "Pretend as if all the refs in 'refs/', along with 'HEAD', are listed. It is equivalent to running 'git log --all'. If set to 'true', the 'from' attribute will be ignored.",
				MarkdownDescription: "Pretend as if all the refs in `refs/`, along with `HEAD`, are listed. It is equivalent to running `git log --all`. If set to `true`, the `from` attribute will be ignored.",
				Computed:            true,
				Optional:            true,
			},
			"since": schema.StringAttribute{
				Description:         "Show commits more recent than a specific date. Date must be in RFC 3339 format, e.g. by using the built-in timestamp/timeadd functions.",
				MarkdownDescription: "Show commits more recent than a specific date. Date must be in RFC 3339 format, e.g. by using the built-in [timestamp](https://www.terraform.io/language/functions/timestamp)/[timeadd](https://www.terraform.io/language/functions/timeadd) functions.",
				Computed:            true,
				Optional:            true,
			},
			"until": schema.StringAttribute{
				Description:         "Show commits older than a specific date. Date must be in RFC 3339 format, e.g. by using the built-in timestamp/timeadd functions.",
				MarkdownDescription: "Show commits older than a specific date. Date must be in RFC 3339 format, e.g. by using the built-in [timestamp](https://www.terraform.io/language/functions/timestamp)/[timeadd](https://www.terraform.io/language/functions/timeadd) functions.",
				Computed:            true,
				Optional:            true,
			},
			"max_count": schema.Int64Attribute{
				Description:         "Limit the number of commits to output.",
				MarkdownDescription: "Limit the number of commits to output.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"skip": schema.Int64Attribute{
				Description:         "Skip first number of commits in output.",
				MarkdownDescription: "Skip first number of commits in output.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"filter_paths": schema.ListAttribute{
				Description:         "Show only commits that are enough to explain how the files that match the specified paths came to be. Note that these are not Git 'pathspec' but rather Go path matchers thus you have to add '/*' for directories yourself.",
				MarkdownDescription: "Show only commits that are enough to explain how the files that match the specified paths came to be. Note that these are not Git `pathspec` but rather Go [path matchers](https://pkg.go.dev/path#Match) thus you have to add `/*` for directories yourself.",
				ElementType:         types.StringType,
				Computed:            true,
				Optional:            true,
			},
			"commits": schema.ListAttribute{
				Description:         "The resulting commit SHA1 hashes ordered as specified by the 'order' attribute.",
				MarkdownDescription: "The resulting commit SHA1 hashes ordered as specified by the `order` attribute.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *LogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_log")

	var inputs logDataSourceModel
	var state logDataSourceModel

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	logOptions := createLogOptions(ctx, repository, &inputs, &resp.Diagnostics)
	if logOptions == nil {
		return
	}

	commits, err := repository.Log(logOptions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read log",
			"Could read log of ["+directory+"] because of: "+err.Error(),
		)
		return
	}
	var hashes []string
	err = commits.ForEach(func(c *object.Commit) error {
		if !inputs.MaxCount.IsNull() && !inputs.MaxCount.IsUnknown() {
			if !inputs.Skip.IsNull() && !inputs.Skip.IsUnknown() {
				if int64(len(hashes)) < inputs.MaxCount.ValueInt64()+inputs.Skip.ValueInt64() {
					hashes = append(hashes, c.Hash.String())
				}
			} else {
				if int64(len(hashes)) < inputs.MaxCount.ValueInt64() {
					hashes = append(hashes, c.Hash.String())
				}
			}
		} else {
			hashes = append(hashes, c.Hash.String())
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read commits",
			"Could not read commits because of: "+err.Error(),
		)
		return
	}
	if !inputs.Skip.IsNull() && !inputs.Skip.IsUnknown() {
		if int64(len(hashes)) >= inputs.Skip.ValueInt64() {
			hashes = hashes[inputs.Skip.ValueInt64():]
		}
	}

	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.All = inputs.All
	state.From = inputs.From
	state.Since = inputs.Since
	state.Until = inputs.Until
	state.MaxCount = inputs.MaxCount
	state.Skip = inputs.Skip
	state.Order = inputs.Order
	state.FilterPaths = inputs.FilterPaths
	state.Commits, _ = types.ListValueFrom(ctx, types.StringType, hashes)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
