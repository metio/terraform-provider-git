/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type dataSourceGitStatusType struct{}

type dataSourceGitStatus struct {
	p gitProvider
}

type dataSourceGitStatusSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	File      types.String `tfsdk:"file"`
	Staging   types.String `tfsdk:"staging"`
	Worktree  types.String `tfsdk:"worktree"`
}

func (r *dataSourceGitStatusType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Fetches the status of a single files in a Git repository.",
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
				MarkdownDescription: "`DEPRECATED`: Only added in order to use the sdkv2 test framework. The file to get status information about.",
				Type:                types.StringType,
				Computed:            true,
			},
			"file": {
				Description: "The file to get status information about.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"staging": {
				Description: "The status of the file in the staging area.",
				Type:        types.StringType,
				Computed:    true,
			},
			"worktree": {
				Description: "The status of the file in the worktree",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func (r *dataSourceGitStatusType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return &dataSourceGitStatus{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitStatus) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	tflog.Debug(ctx, "Reading Git file status")

	var inputs dataSourceGitStatusSchema
	var state dataSourceGitStatusSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	fileName := inputs.File.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	state.Directory = types.String{Value: directory}
	state.Id = types.String{Value: fileName}
	state.File = types.String{Value: fileName}

	worktree, err := repository.Worktree()
	if err == git.ErrIsBareRepository {
		tflog.Trace(ctx, "read worktree of bare repository", map[string]interface{}{
			"directory": directory,
		})
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read worktree",
			"Could not read worktree because of: "+err.Error(),
		)
		return
	} else {
		tflog.Trace(ctx, "read worktree", map[string]interface{}{
			"directory": directory,
		})

		status, err := worktree.Status()
		if err != nil {
			resp.Diagnostics.AddError(
				"Cannot read status",
				"Could not read status because of: "+err.Error(),
			)
			return
		}
		tflog.Trace(ctx, "read status", map[string]interface{}{
			"directory": directory,
			"status":    status.String(),
		})

		fileStatus := status.File(fileName)
		tflog.Trace(ctx, "read file status", map[string]interface{}{
			"directory": directory,
			"file":      fileName,
			"staging":   fileStatus.Staging,
			"worktree":  fileStatus.Worktree,
		})
		state.Staging = types.String{Value: string(fileStatus.Staging)}
		state.Worktree = types.String{Value: string(fileStatus.Worktree)}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}