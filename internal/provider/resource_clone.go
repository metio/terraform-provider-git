/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
)

type CloneResource struct{}

var (
	_ resource.Resource               = (*CloneResource)(nil)
	_ resource.ResourceWithModifyPlan = (*CloneResource)(nil)
)

type cloneResourceModel struct {
	Directory     types.String `tfsdk:"directory"`
	Id            types.String `tfsdk:"id"`
	Bare          types.Bool   `tfsdk:"bare"`
	RemoteName    types.String `tfsdk:"remote_name"`
	ReferenceName types.String `tfsdk:"reference_name"`
	URL           types.String `tfsdk:"url"`
	Auth          types.Object `tfsdk:"auth"`
	SHA1          types.String `tfsdk:"sha1"`
}

func NewCloneResource() resource.Resource {
	return &CloneResource{}
}

func (r *CloneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clone"
}

func (r *CloneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Clones a Git repository similar to 'git clone'.",
		MarkdownDescription: "Clones a Git repository similar to `git clone`.",
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
				Description:         "The import ID to import this resource which is equal to the value of the 'directory' attribute.",
				MarkdownDescription: "The import ID to import this resource which is equal to the value of the `directory` attribute.",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				Description:         "The (possibly remote) repository URL to clone from.",
				MarkdownDescription: "The (possibly remote) repository URL to clone from.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bare": schema.BoolAttribute{
				Description:         "Whether we should perform a bare clone. Defaults to 'false'.",
				MarkdownDescription: "Whether we should perform a bare clone. Defaults to `false`.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.DefaultBool(false),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"remote_name": schema.StringAttribute{
				Description:         "Name of the remote to be added. Defaults to 'origin'.",
				MarkdownDescription: "Name of the remote to be added. Defaults to 'origin'.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.DefaultString("origin"),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"reference_name": schema.StringAttribute{
				Description:         "Name of the remote to be added. Defaults to 'main'.",
				MarkdownDescription: "Name of the remote to be added. Defaults to 'main'.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.DefaultString("main"),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auth": schema.SingleNestedAttribute{
				Description:         "The authentication credentials, if required, to use with the remote repository.",
				MarkdownDescription: "The authentication credentials, if required, to use with the remote repository.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"basic": schema.SingleNestedAttribute{
						Description:         "Configure basic auth authentication.",
						MarkdownDescription: "Configure basic auth authentication.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								Description:         "The basic auth username.",
								MarkdownDescription: "The basic auth username.",
								Required:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"password": schema.StringAttribute{
								Description:         "The basic auth password.",
								MarkdownDescription: "The basic auth password.",
								Required:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_key")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_password")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_agent")),
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.RequiresReplace(),
						},
					},
					"bearer": schema.StringAttribute{
						Description:         "Configure HTTP bearer token authentication. Note that services like GitHub use basic auth with your OAuth2 personal access token as the password.",
						MarkdownDescription: "Configure HTTP bearer token authentication. **Note**: Services like GitHub use basic auth with your OAuth2 personal access token as the password.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic")),
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_key")),
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_password")),
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_agent")),
							stringvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"ssh_key": schema.SingleNestedAttribute{
						Description:         "Configure SSH public/private key authentication.",
						MarkdownDescription: "Configure SSH public/private key authentication.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								Description:         "The SSH auth username.",
								MarkdownDescription: "The SSH auth username.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									modifiers.DefaultString("git"),
									stringplanmodifier.RequiresReplace(),
								},
							},
							"password": schema.StringAttribute{
								Description:         "The SSH key password.",
								MarkdownDescription: "The SSH key password.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									modifiers.DefaultString(""),
									stringplanmodifier.RequiresReplace(),
								},
							},
							"private_key_path": schema.StringAttribute{
								Description:         "The absolute path to the private SSH key.",
								MarkdownDescription: "The absolute path to the private SSH key.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("private_key_pem")),
									stringvalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("private_key_pem")),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"private_key_pem": schema.StringAttribute{
								Description:         "The private SSH key in PEM format.",
								MarkdownDescription: "The private SSH key in PEM format.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("private_key_path")),
									stringvalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("private_key_path")),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"known_hosts": schema.SetAttribute{
								Description:         "The list of known hosts files to accept. If none are specified, system defaults will be used.",
								MarkdownDescription: "The list of known hosts files to accept. If none are specified, system defaults will be used.",
								ElementType:         types.StringType,
								Optional:            true,
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.RequiresReplace(),
								},
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_password")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_agent")),
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
					},
					"ssh_password": schema.SingleNestedAttribute{
						Description:         "Configure password based SSH authentication.",
						MarkdownDescription: "Configure password based SSH authentication.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								Description:         "The SSH username.",
								MarkdownDescription: "The SSH username.",
								Required:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"password": schema.StringAttribute{
								Description:         "The SSH password.",
								MarkdownDescription: "The SSH password.",
								Required:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"known_hosts": schema.SetAttribute{
								Description:         "The list of known hosts files to accept. If none are specified, system defaults will be used.",
								MarkdownDescription: "The list of known hosts files to accept. If none are specified, system defaults will be used.",
								ElementType:         types.StringType,
								Optional:            true,
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.RequiresReplace(),
								},
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_key")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_agent")),
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.RequiresReplace(),
						},
					},
					"ssh_agent": schema.SingleNestedAttribute{
						Description:         "Configure SSH agent based authentication.",
						MarkdownDescription: "Configure SSH agent based authentication.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								Description:         "The system username of the user talking to the SSH agent. Use an empty string in order to automatically fetch this.",
								MarkdownDescription: "The system username of the user talking to the SSH agent. Use an empty string in order to automatically fetch this.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									modifiers.DefaultString(""),
									stringplanmodifier.RequiresReplace(),
								},
							},
							"known_hosts": schema.SetAttribute{
								Description:         "The list of known hosts files to accept. If none are specified, system defaults will be used.",
								MarkdownDescription: "The list of known hosts files to accept. If none are specified, system defaults will be used.",
								ElementType:         types.StringType,
								Optional:            true,
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.RequiresReplace(),
								},
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_key")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_password")),
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.RequiresReplace(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"sha1": schema.StringAttribute{
				Description:         "The SHA1 hash of the created commit.",
				MarkdownDescription: "The SHA1 hash of the created commit.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *CloneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create resource git_clone")

	var inputs cloneResourceModel
	var state cloneResourceModel

	diags := req.Plan.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()
	bare := inputs.Bare.ValueBool()

	options := CreateCloneOptions(ctx, &inputs, &resp.Diagnostics)
	if options == nil {
		return
	}

	repository, err := git.PlainCloneContext(ctx, directory, bare, options)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot clone repository",
			"Could not clone repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "cloned repository", map[string]interface{}{
		"directory": directory,
		"bare":      bare,
	})

	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.Bare = inputs.Bare
	state.URL = inputs.URL
	state.RemoteName = inputs.RemoteName
	state.ReferenceName = inputs.ReferenceName
	state.Auth = inputs.Auth
	state.SHA1 = types.StringNull()

	head, err := repository.Head()
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot read HEAD",
			"Could not read HEAD of repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}
	if head != nil {
		state.SHA1 = types.StringValue(head.Hash().String())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CloneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read resource git_clone")
	// NO-OP: All data is already in Terraform state
}

func (r *CloneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update resource git_clone")
	// NO-OP: All data is already in Terraform state
}

func (r *CloneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete resource git_clone")

	var state cloneResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.ValueString()
	bare := state.Bare.ValueBool()

	if !bare {
		repository := openRepository(ctx, directory, &resp.Diagnostics)

		if repository != nil && repository.Storer != nil {
			storage, ok := repository.Storer.(*filesystem.Storage)

			if ok {
				err := os.RemoveAll(storage.Filesystem().Root())
				if err != nil {
					resp.Diagnostics.AddError(
						"Cannot delete repository",
						"Could not delete git repository ["+directory+"] because of: "+err.Error(),
					)
					return
				}
			}
		}
	}
}

func (r *CloneResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	tflog.Debug(ctx, "ModifyPlan resource git_clone")

	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to modify it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to modify it
		return
	}

	var inputs cloneResourceModel
	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()
	url := inputs.URL.ValueString()

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	head, err := repository.Head()
	if err != nil {
		diags.AddError(
			"Cannot read HEAD",
			"Could not read HEAD of git repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	localHeadHash := head.Hash()

	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{strings.ReplaceAll(url, "/", "\\")},
	})
	refs, err := remote.List(&git.ListOptions{
		PeelingOption: git.AppendPeeled,
		Auth:          authOptions(ctx, inputs.Auth, &diags),
	})
	if err != nil {
		diags.AddError(
			"Cannot list remote",
			"Could not list remote of git repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	for _, ref := range refs {
		if ref.Name().IsBranch() {
			var expectedRefName string
			if inputs.ReferenceName.ValueString() == "" {
				expectedRefName = "refs/heads/main"
			} else {
				expectedRefName = fmt.Sprintf("refs/heads/%s", inputs.ReferenceName.ValueString())
			}
			if expectedRefName == ref.Name().String() {
				if localHeadHash.String() != ref.Hash().String() {
					sha1 := path.Root("sha1")
					resp.Plan.SetAttribute(ctx, sha1, ref.Hash().String())
					resp.RequiresReplace = append(resp.RequiresReplace, sha1)
				}
			}
		}
	}
}
