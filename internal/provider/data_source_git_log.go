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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
)

type dataSourceGitLogType struct{}

type dataSourceGitLog struct {
	p gitProvider
}

type dataSourceGitLogSchema struct {
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

func (r *dataSourceGitLogType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Fetches the commit log of a Git repository similar to 'git log'.",
		MarkdownDescription: "Fetches the commit log of a Git repository similar to `git log`.",
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
				Description:         "The same value as the 'directory' attribute.",
				MarkdownDescription: "The same value as the `directory` attribute.",
				Type:                types.StringType,
				Computed:            true,
			},
			"from": {
				Description:         "When set the log will only contain commits reachable from it. If this option is not set, 'HEAD' will be used as the default. Can be any revision that 'go-git' supports. See https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision for details.",
				MarkdownDescription: "When set the log will only contain commits reachable from it. If this option is not set, `HEAD` will be used as the default. Can be any [revision](https://www.git-scm.com/docs/gitrevisions) that `go-git` [supports](https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision).",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"order": {
				Description:         "The traversal algorithm to use while listing commits. Defaults to 'time' which is similar to 'git log'. Other values are 'depth' and 'breadth' for depth- or breadth-first traversal.",
				MarkdownDescription: "The traversal algorithm to use while listing commits. Defaults to `time` which is similar to `git log`. Other values are `depth` and `breadth` for depth- or breadth-first traversal.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.String{Value: "time"}),
				},
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf(
						"time",
						"depth",
						"breadth",
					),
				},
			},
			"all": {
				Description:         "Pretend as if all the refs in 'refs/', along with 'HEAD', are listed. It is equivalent to running 'git log --all'. If set to 'true', the 'from' attribute will be ignored.",
				MarkdownDescription: "Pretend as if all the refs in `refs/`, along with `HEAD`, are listed. It is equivalent to running `git log --all`. If set to `true`, the `from` attribute will be ignored.",
				Type:                types.BoolType,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.Bool{Value: false}),
				},
			},
			"since": {
				Description:         "Show commits more recent than a specific date. Date must be in RFC 3339 format, e.g. by using the built-in timestamp/timeadd functions.",
				MarkdownDescription: "Show commits more recent than a specific date. Date must be in RFC 3339 format, e.g. by using the built-in [timestamp](https://www.terraform.io/language/functions/timestamp)/[timeadd](https://www.terraform.io/language/functions/timeadd) functions.",
				Type:                types.StringType,
				Computed:            true,
				Optional:            true,
			},
			"until": {
				Description:         "Show commits older than a specific date. Date must be in RFC 3339 format, e.g. by using the built-in timestamp/timeadd functions.",
				MarkdownDescription: "Show commits older than a specific date. Date must be in RFC 3339 format, e.g. by using the built-in [timestamp](https://www.terraform.io/language/functions/timestamp)/[timeadd](https://www.terraform.io/language/functions/timeadd) functions.",
				Type:                types.StringType,
				Computed:            true,
				Optional:            true,
			},
			"max_count": {
				Description:         "Limit the number of commits to output.",
				MarkdownDescription: "Limit the number of commits to output.",
				Type:                types.Int64Type,
				Computed:            true,
				Optional:            true,
				Validators: []tfsdk.AttributeValidator{
					int64validator.AtLeast(0),
				},
			},
			"skip": {
				Description:         "Skip first number of commits in output.",
				MarkdownDescription: "Skip first number of commits in output.",
				Type:                types.Int64Type,
				Computed:            true,
				Optional:            true,
				Validators: []tfsdk.AttributeValidator{
					int64validator.AtLeast(0),
				},
			},
			"filter_paths": {
				Description:         "Show only commits that are enough to explain how the files that match the specified paths came to be. Note that these are not Git 'pathspec' but rather Go path matchers thus you have to add '/*' for directories yourself.",
				MarkdownDescription: "Show only commits that are enough to explain how the files that match the specified paths came to be. Note that these are not Git `pathspec` but rather Go [path matchers](https://pkg.go.dev/path#Match) thus you have to add `/*` for directories yourself.",
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Computed: true,
				Optional: true,
			},
			"commits": {
				Description:         "The resulting commit SHA1 hashes ordered as specified by the 'order' attribute.",
				MarkdownDescription: "The resulting commit SHA1 hashes ordered as specified by the `order` attribute.",
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Computed: true,
			},
		},
	}, nil
}

func (r *dataSourceGitLogType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return &dataSourceGitLog{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitLog) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data source git_log")

	var inputs dataSourceGitLogSchema
	var state dataSourceGitLogSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value

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
				if int64(len(hashes)) < inputs.MaxCount.Value+inputs.Skip.Value {
					hashes = append(hashes, c.Hash.String())
				}
			} else {
				if int64(len(hashes)) < inputs.MaxCount.Value {
					hashes = append(hashes, c.Hash.String())
				}
			}
		} else {
			hashes = append(hashes, c.Hash.String())
		}
		return nil
	})
	if !inputs.Skip.IsNull() && !inputs.Skip.IsUnknown() {
		if int64(len(hashes)) >= inputs.Skip.Value {
			hashes = hashes[inputs.Skip.Value:]
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
	state.Commits = stringsToList(hashes)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
