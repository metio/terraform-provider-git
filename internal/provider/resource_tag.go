/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
	"strings"
)

type TagResource struct{}

var (
	_ resource.Resource = (*TagResource)(nil)
)

type tagResourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Message   types.String `tfsdk:"message"`
	Revision  types.String `tfsdk:"revision"`
	SHA1      types.String `tfsdk:"sha1"`
}

func NewTagResource() resource.Resource {
	return &TagResource{}
}

func (r *TagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *TagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage Git tags similar to 'git tag'.",
		MarkdownDescription: "Manage Git tags similar to `git tag`.",
		Attributes: map[string]schema.Attribute{
			"directory": schema.StringAttribute{
				Description:         "The path to the local Git repository.",
				MarkdownDescription: "The path to the local Git repository.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description:         "The import ID to import this resource which has the form 'directory|name'",
				MarkdownDescription: "The import ID to import this resource which has the form `'directory|name'`",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "The name of the Git tag to add.",
				MarkdownDescription: "The name of the Git tag to add.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"revision": schema.StringAttribute{
				Description:         "The revision of the commit to tag. Can be any value that 'go-git' supports. If none is specified, 'HEAD' will be tagged.",
				MarkdownDescription: "The [revision](https://www.git-scm.com/docs/gitrevisions) of the commit to tag. Can be any value that `go-git` [supports](https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision). If none is specified, `HEAD` will be tagged.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					modifiers.DefaultString("HEAD"),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sha1": schema.StringAttribute{
				Description:         "The SHA1 hash of the resolved revision.",
				MarkdownDescription: "The SHA1 hash of the resolved revision.",
				Computed:            true,
			},
			"message": schema.StringAttribute{
				Description:         "The tag message to use. Note that by specifying a message, an annotated tag will be created.",
				MarkdownDescription: "The tag message to use. Note that by specifying a message, an annotated tag will be created.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *TagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create resource git_tag")

	var inputs tagResourceModel
	diags := req.Plan.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()
	tagName := inputs.Name.ValueString()

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	hash := resolveRevision(ctx, repository, inputs.Revision.ValueString(), &resp.Diagnostics)
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

	var state tagResourceModel
	state.Directory = inputs.Directory
	state.Id = types.StringValue(fmt.Sprintf("%s|%s", directory, tagName))
	state.Name = inputs.Name
	state.Message = inputs.Message
	state.Revision = inputs.Revision
	state.SHA1 = types.StringValue(hash.String())

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read resource git_tag")

	var state tagResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.ValueString()
	tagName := state.Name.ValueString()

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

	var newState tagResourceModel
	newState.Directory = state.Directory
	newState.Id = types.StringValue(fmt.Sprintf("%s|%s", directory, tagName))
	newState.Name = state.Name
	newState.Revision = state.Revision
	newState.SHA1 = types.StringValue(tagReference.Hash().String())
	if tagObject == nil {
		newState.Message = types.StringNull()
	} else {
		newState.Message = types.StringValue(strings.TrimSpace(tagObject.Message))
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TagResource) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update resource git_tag")
	// NO-OP: all attributes require replace, thus Delete and Create methods will be called
}

func (r *TagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete resource git_tag")

	var state tagResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.ValueString()
	tagName := state.Name.ValueString()

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

func (r *TagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

	var state tagResourceModel
	state.Directory = types.StringValue(directory)
	state.Id = types.StringValue(fmt.Sprintf("%s|%s", directory, tagName))
	state.Name = types.StringValue(tagName)
	state.Revision = types.StringValue(revision)
	state.SHA1 = types.StringValue(tagReference.Hash().String())
	if tagObject == nil {
		state.Message = types.StringNull()
	} else {
		state.Message = types.StringValue(strings.TrimSpace(tagObject.Message))
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
