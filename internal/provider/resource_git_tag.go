/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
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
	SHA1      types.String `tfsdk:"sha1"`
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
				MarkdownDescription: "`DEPRECATED`: Only added in order to use the sdkv2 test framework. The name of the Git tag to add.",
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
			"sha1": {
				Description: "The SHA1 checksum of the commit to tag.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"message": {
				Description: "The tag message to use. Note that by specifying a message, an annotated tag will be created.",
				Type:        types.StringType,
				Optional:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
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

	var config resourceGitTagSchema
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := config.Directory.Value
	tagName := config.Name.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	reference, err := createTagReference(repository, config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot create tag reference",
			"Could not create tag reference in git repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created tag reference", map[string]interface{}{
		"directory": directory,
		"tag":       tagName,
		"reference": reference.Hash(),
	})

	_, err = repository.CreateTag(tagName, reference.Hash(), createOptions(config))
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot create tag",
			"Could not create tag ["+tagName+"] in git repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created tag", map[string]interface{}{
		"directory": directory,
		"tag":       tagName,
	})

	var state resourceGitTagSchema
	state.Directory = config.Directory
	state.Id = config.Name
	state.Name = config.Name
	state.Message = config.Message
	state.SHA1 = types.String{Value: reference.Hash().String()}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitTag) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Debug(ctx, "Reading Git tag")

	var state resourceGitTagSchema
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.Value
	tagName := state.Name.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	tagReference := getTagReference(ctx, repository, tagName, &resp.Diagnostics)
	if tagReference == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var newState resourceGitTagSchema
	newState.Directory = state.Directory
	newState.Id = state.Name
	newState.Name = state.Name
	newState.SHA1 = types.String{Value: tagReference.Hash().String()}
	tag, err := repository.TagObject(tagReference.Hash())
	if err == plumbing.ErrObjectNotFound {
		newState.Message = types.String{Null: true}
	} else {
		newState.Message = types.String{Value: strings.TrimSpace(tag.Message)}
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitTag) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "Updating Git tag")
	// NO-OP: all attributes require replace, thus Delete and Create methods will be called
}

func (r *resourceGitTag) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "Removing Git tag")

	var state resourceGitTagSchema
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.Value
	tagName := state.Name.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	err := repository.DeleteTag(tagName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot delete tag",
			"Could not delete tag ["+tagName+"] in git repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}
}

func (r *resourceGitTag) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	id := req.ID
	idParts := strings.Split(id, "|")

	if len(idParts) < 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected import identifier",
			fmt.Sprintf("Expected import identifier with format: 'directory|tag-name|sha1|message' Got: %q", id),
		)
		return
	}

	var state resourceGitTagSchema

	state.Directory = types.String{Value: idParts[0]}
	state.Id = types.String{Value: idParts[1]}
	state.Name = types.String{Value: idParts[1]}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
