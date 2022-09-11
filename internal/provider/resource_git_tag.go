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
	"github.com/metio/terraform-provider-git/internal/modifiers"
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
	Revision  types.String `tfsdk:"revision"`
	SHA1      types.String `tfsdk:"sha1"`
}

func (r *resourceGitTagType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Manage Git tags similar to 'git tag'.",
		MarkdownDescription: "Manage Git tags similar to `git tag`.",
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description:         "The path to the local Git repository.",
				MarkdownDescription: "The path to the local Git repository.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"id": {
				Description:         "The import ID to import this resource which has the form 'directory|name'",
				MarkdownDescription: "The import ID to import this resource which has the form `'directory|name'`",
				Type:                types.StringType,
				Computed:            true,
			},
			"name": {
				Description:         "The name of the Git tag to add.",
				MarkdownDescription: "The name of the Git tag to add.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"revision": {
				Description:         "The revision of the commit to tag. Can be any value that 'go-git' supports. If none is specified, 'HEAD' will be tagged.",
				MarkdownDescription: "The [revision](https://www.git-scm.com/docs/gitrevisions) of the commit to tag. Can be any value that `go-git` [supports](https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision). If none is specified, `HEAD` will be tagged.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.String{Value: "HEAD"}),
					resource.RequiresReplace(),
				},
			},
			"sha1": {
				Description:         "The SHA1 hash of the resolved revision.",
				MarkdownDescription: "The SHA1 hash of the resolved revision.",
				Type:                types.StringType,
				Computed:            true,
			},
			"message": {
				Description:         "The tag message to use. Note that by specifying a message, an annotated tag will be created.",
				MarkdownDescription: "The tag message to use. Note that by specifying a message, an annotated tag will be created.",
				Type:                types.StringType,
				Optional:            true,
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
	tflog.Debug(ctx, "Create resource git_tag")

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

	// NOTE: It seems default values are not working?
	if inputs.Revision.IsNull() {
		inputs.Revision = types.String{Value: "HEAD"}
	}

	hash := resolveRevision(ctx, repository, inputs.Revision.Value, &resp.Diagnostics)
	if hash == nil {
		return
	}

	_, err := repository.CreateTag(tagName, *hash, createTagOptions(inputs))
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
	state.Id = types.String{Value: fmt.Sprintf("%s|%s", directory, tagName)}
	state.Name = inputs.Name
	state.Message = inputs.Message
	state.Revision = inputs.Revision
	state.SHA1 = types.String{Value: hash.String()}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitTag) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read resource git_tag")

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
	newState.Id = types.String{Value: fmt.Sprintf("%s|%s", directory, tagName)}
	newState.Name = state.Name
	newState.Revision = state.Revision
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

func (r *resourceGitTag) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update resource git_tag")
	// NO-OP: all attributes require replace, thus Delete and Create methods will be called
}

func (r *resourceGitTag) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete resource git_tag")

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
	tflog.Debug(ctx, "ImportState resource git_tag")

	id := req.ID
	idParts := strings.Split(id, "|")

	if len(idParts) < 2 || len(idParts) > 3 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected import identifier",
			fmt.Sprintf("Expected import identifier with format: 'path/to/your/git/repository|name-of-your-tag|revision' Got: %q", id),
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

	var revision string
	if len(idParts) == 2 {
		revision = "HEAD"
	} else {
		revision = idParts[2]
	}

	var state resourceGitTagSchema
	state.Directory = types.String{Value: directory}
	state.Id = types.String{Value: id}
	state.Name = types.String{Value: tagName}
	state.Revision = types.String{Value: revision}
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
