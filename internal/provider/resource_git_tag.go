/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
		MarkdownDescription: "Manage Git tags with `git tag`.",
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description: "The path to the local Git repository.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
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
					resource.RequiresReplace(),
				},
			},
			"sha1": {
				MarkdownDescription: "The SHA1 checksum of the commit to tag. If none is specified, `HEAD` will be tagged.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"message": {
				Description: "The tag message to use. Note that by specifying a message, an annotated tag will be created.",
				Type:        types.StringType,
				Optional:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (r *resourceGitTagType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return &resourceGitTag{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *resourceGitTag) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Git tag")

	var inputs resourceGitTagSchema
	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value
	tagName := inputs.Name.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

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
		"tag":       tagName,
		"reference": reference.Hash(),
	})

	_, err = repository.CreateTag(tagName, reference.Hash(), createTagOptions(inputs))
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
	state.Directory = inputs.Directory
	state.Id = inputs.Name
	state.Name = inputs.Name
	state.Message = inputs.Message
	state.SHA1 = types.String{Value: reference.Hash().String()}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitTag) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	tagObject, err := getTagObject(ctx, repository, tagReference.Hash(), &resp.Diagnostics)
	if err != nil {
		return
	}

	var newState resourceGitTagSchema
	newState.Directory = state.Directory
	newState.Id = state.Name
	newState.Name = state.Name
	newState.SHA1 = types.String{Value: tagReference.Hash().String()}
	if tagObject == nil {
		newState.Message = types.String{Null: true}
	} else {
		newState.Message = types.String{Value: strings.TrimSpace(tagObject.Message)}
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitTag) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Git tag")
	// NO-OP: all attributes require replace, thus Delete and Create methods will be called
}

func (r *resourceGitTag) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

func (r *resourceGitTag) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	idParts := strings.Split(id, "|")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected import identifier",
			fmt.Sprintf("Expected import identifier with format: 'path/to/your/git/repository|name-of-your-tag' Got: %q", id),
		)
		return
	}

	directory := idParts[0]
	tagName := idParts[1]

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	tagReference := getTagReference(ctx, repository, tagName, &resp.Diagnostics)
	if tagReference == nil {
		return
	}

	tagObject, err := getTagObject(ctx, repository, tagReference.Hash(), &resp.Diagnostics)
	if err != nil {
		return
	}

	var state resourceGitTagSchema
	state.Directory = types.String{Value: directory}
	state.Id = types.String{Value: tagName}
	state.Name = types.String{Value: tagName}
	state.SHA1 = types.String{Value: tagReference.Hash().String()}
	if tagObject == nil {
		state.Message = types.String{Null: true}
	} else {
		state.Message = types.String{Value: strings.TrimSpace(tagObject.Message)}
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
