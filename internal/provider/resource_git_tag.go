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
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type resourceGitTagType struct{}

type resourceGitTag struct {
	p gitProvider
}

type resourceGitTagSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Message   types.String `tfsdk:"message"`
	SHA1      types.String `tfsdk:"commit_sha1"`
}

func (c *resourceGitTagType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Adds a new Git tag to a repository.",
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description: "The path to the local Git repository.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"id": {
				MarkdownDescription: "`DEPRECATED`: Only added in order to use the sdkv2 test framework. The path to the local Git repository.",
				Type:                types.StringType,
				Computed:            true,
			},
			"name": {
				Description: "The name of the Git tag to add.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"message": {
				Description: "The tag message to use. Note that by specifying a message, an annotated tag will be created.",
				Type:        types.StringType,
				Optional:    true,
			},
			"commit_sha1": {
				Description: "The SHA1 checksum of the commit to tag.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
		},
	}, nil
}

func (r *resourceGitTagType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return &resourceGitTag{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *resourceGitTag) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, "Creating Git tag")

	var inputs resourceGitTagSchema
	var output resourceGitTagSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	name := inputs.Name.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	tflog.Trace(ctx, "opened repository", map[string]interface{}{
		"directory": directory,
	})

	reference, err := createTagReference(repository, inputs)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot create tag reference",
			"Could not create tag reference in git repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created tag reference", map[string]interface{}{
		"directory": directory,
		"tag":       name,
		"reference": reference.Hash(),
	})

	tag, err := repository.CreateTag(name, reference.Hash(), createOptions(inputs))
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot create tag",
			"Could not create tag ["+name+"] in git repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created tag", map[string]interface{}{
		"directory": directory,
		"tag":       tag.Name().Short(),
	})

	output.Directory = types.String{Value: directory}
	output.Id = types.String{Value: directory}
	output.Name = types.String{Value: tag.Name().Short()}
	output.Message = inputs.Message
	output.SHA1 = types.String{Value: reference.Hash().String()}

	diags = resp.State.Set(ctx, &output)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func createTagReference(repository *git.Repository, inputs resourceGitTagSchema) (*plumbing.Reference, error) {
	if inputs.SHA1.IsNull() || inputs.SHA1.IsUnknown() {
		head, err := repository.Head()
		if err != nil {
			return nil, err
		}
		return head, nil
	}

	return plumbing.NewHashReference("tag", plumbing.NewHash(inputs.SHA1.Value)), nil
}

func createOptions(inputs resourceGitTagSchema) *git.CreateTagOptions {
	if inputs.Message.IsNull() || inputs.Message.IsUnknown() {
		return nil
	}
	return &git.CreateTagOptions{
		Message: inputs.Message.Value,
	}
}

func (r *resourceGitTag) Read(ctx context.Context, _ tfsdk.ReadResourceRequest, _ *tfsdk.ReadResourceResponse) {
	// NO-OP: all there is to read is in the State, and response is already populated with that.
	tflog.Debug(ctx, "Reading Git tag from state")
}

func (r *resourceGitTag) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "Updating Git tag")
	updatedUsingPlan(ctx, &req, resp, &resourceGitTagSchema{})
}

func (r *resourceGitTag) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "Removing Git tag")

	var inputs resourceGitTagSchema

	diags := req.State.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	name := inputs.Name.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	err := repository.DeleteTag(name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot delete tag",
			"Could not delete tag ["+name+"] in git repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}
}

func (r *resourceGitTag) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
