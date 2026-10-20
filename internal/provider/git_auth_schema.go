/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/metio/terraform-provider-git/internal/modifiers"
)

func authResourceSchema() resourceschema.SingleNestedAttribute {
	return resourceschema.SingleNestedAttribute{
		Description:         "The authentication credentials, if required, to use with the remote repository.",
		MarkdownDescription: "The authentication credentials, if required, to use with the remote repository.",
		Optional:            true,
		Attributes: map[string]resourceschema.Attribute{
			"basic": resourceschema.SingleNestedAttribute{
				Description:         "Configure basic auth authentication.",
				MarkdownDescription: "Configure basic auth authentication.",
				Optional:            true,
				Attributes: map[string]resourceschema.Attribute{
					"username": resourceschema.StringAttribute{
						Description:         "The basic auth username.",
						MarkdownDescription: "The basic auth username.",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"password": resourceschema.StringAttribute{
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
			"bearer": resourceschema.StringAttribute{
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
			"ssh_key": resourceschema.SingleNestedAttribute{
				Description:         "Configure SSH public/private key authentication.",
				MarkdownDescription: "Configure SSH public/private key authentication.",
				Optional:            true,
				Attributes: map[string]resourceschema.Attribute{
					"username": resourceschema.StringAttribute{
						Description:         "The SSH auth username.",
						MarkdownDescription: "The SSH auth username.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString("git"),
							stringplanmodifier.RequiresReplace(),
						},
					},
					"password": resourceschema.StringAttribute{
						Description:         "The SSH key password.",
						MarkdownDescription: "The SSH key password.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString(""),
							stringplanmodifier.RequiresReplace(),
						},
					},
					"private_key_path": resourceschema.StringAttribute{
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
					"private_key_pem": resourceschema.StringAttribute{
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
					"known_hosts": resourceschema.SetAttribute{
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
			"ssh_password": resourceschema.SingleNestedAttribute{
				Description:         "Configure password based SSH authentication.",
				MarkdownDescription: "Configure password based SSH authentication.",
				Optional:            true,
				Attributes: map[string]resourceschema.Attribute{
					"username": resourceschema.StringAttribute{
						Description:         "The SSH username.",
						MarkdownDescription: "The SSH username.",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"password": resourceschema.StringAttribute{
						Description:         "The SSH password.",
						MarkdownDescription: "The SSH password.",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"known_hosts": resourceschema.SetAttribute{
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
			"ssh_agent": resourceschema.SingleNestedAttribute{
				Description:         "Configure SSH agent based authentication.",
				MarkdownDescription: "Configure SSH agent based authentication.",
				Optional:            true,
				Attributes: map[string]resourceschema.Attribute{
					"username": resourceschema.StringAttribute{
						Description:         "The system username of the user talking to the SSH agent. Use an empty string in order to automatically fetch this.",
						MarkdownDescription: "The system username of the user talking to the SSH agent. Use an empty string in order to automatically fetch this.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.DefaultString(""),
							stringplanmodifier.RequiresReplace(),
						},
					},
					"known_hosts": resourceschema.SetAttribute{
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
	}
}

func authDataSourceSchema() datasourceschema.SingleNestedAttribute {
	return datasourceschema.SingleNestedAttribute{
		Description:         "The authentication credentials, if required, to use with the remote repository.",
		MarkdownDescription: "The authentication credentials, if required, to use with the remote repository.",
		Optional:            true,
		Attributes: map[string]datasourceschema.Attribute{
			"basic": datasourceschema.SingleNestedAttribute{
				Description:         "Configure basic auth authentication.",
				MarkdownDescription: "Configure basic auth authentication.",
				Optional:            true,
				Attributes: map[string]datasourceschema.Attribute{
					"username": datasourceschema.StringAttribute{
						Description:         "The basic auth username.",
						MarkdownDescription: "The basic auth username.",
						Required:            true,
					},
					"password": datasourceschema.StringAttribute{
						Description:         "The basic auth password.",
						MarkdownDescription: "The basic auth password.",
						Required:            true,
						Sensitive:           true,
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
			},
			"bearer": datasourceschema.StringAttribute{
				Description:         "Configure HTTP bearer token authentication. Note that services like GitHub use basic auth with your OAuth2 personal access token as the password.",
				MarkdownDescription: "Configure HTTP bearer token authentication. **Note**: Services like GitHub use basic auth with your OAuth2 personal access token as the password.",
				Optional:            true,
				Sensitive:           true,
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
			},
			"ssh_key": datasourceschema.SingleNestedAttribute{
				Description:         "Configure SSH public/private key authentication.",
				MarkdownDescription: "Configure SSH public/private key authentication.",
				Optional:            true,
				Attributes: map[string]datasourceschema.Attribute{
					"username": datasourceschema.StringAttribute{
						Description:         "The SSH auth username. Defaults to 'git'.",
						MarkdownDescription: "The SSH auth username. Defaults to `git`.",
						Optional:            true,
					},
					"password": datasourceschema.StringAttribute{
						Description:         "The SSH key password.",
						MarkdownDescription: "The SSH key password.",
						Optional:            true,
						Sensitive:           true,
					},
					"private_key_path": datasourceschema.StringAttribute{
						Description:         "The absolute path to the private SSH key.",
						MarkdownDescription: "The absolute path to the private SSH key.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("private_key_pem")),
							stringvalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("private_key_pem")),
						},
					},
					"private_key_pem": datasourceschema.StringAttribute{
						Description:         "The private SSH key in PEM format.",
						MarkdownDescription: "The private SSH key in PEM format.",
						Optional:            true,
						Sensitive:           true,
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("private_key_path")),
							stringvalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("private_key_path")),
						},
					},
					"known_hosts": datasourceschema.SetAttribute{
						Description:         "The list of known hosts files to accept. If none are specified, system defaults will be used.",
						MarkdownDescription: "The list of known hosts files to accept. If none are specified, system defaults will be used.",
						ElementType:         types.StringType,
						Optional:            true,
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
			"ssh_password": datasourceschema.SingleNestedAttribute{
				Description:         "Configure password based SSH authentication.",
				MarkdownDescription: "Configure password based SSH authentication.",
				Optional:            true,
				Attributes: map[string]datasourceschema.Attribute{
					"username": datasourceschema.StringAttribute{
						Description:         "The SSH username.",
						MarkdownDescription: "The SSH username.",
						Required:            true,
					},
					"password": datasourceschema.StringAttribute{
						Description:         "The SSH password.",
						MarkdownDescription: "The SSH password.",
						Required:            true,
						Sensitive:           true,
					},
					"known_hosts": datasourceschema.SetAttribute{
						Description:         "The list of known hosts files to accept. If none are specified, system defaults will be used.",
						MarkdownDescription: "The list of known hosts files to accept. If none are specified, system defaults will be used.",
						ElementType:         types.StringType,
						Optional:            true,
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
			},
			"ssh_agent": datasourceschema.SingleNestedAttribute{
				Description:         "Configure SSH agent based authentication.",
				MarkdownDescription: "Configure SSH agent based authentication.",
				Optional:            true,
				Attributes: map[string]datasourceschema.Attribute{
					"username": datasourceschema.StringAttribute{
						Description:         "The system username of the user talking to the SSH agent. Use an empty string in order to automatically fetch this.",
						MarkdownDescription: "The system username of the user talking to the SSH agent. Use an empty string in order to automatically fetch this.",
						Optional:            true,
					},
					"known_hosts": datasourceschema.SetAttribute{
						Description:         "The list of known hosts files to accept. If none are specified, system defaults will be used.",
						MarkdownDescription: "The list of known hosts files to accept. If none are specified, system defaults will be used.",
						ElementType:         types.StringType,
						Optional:            true,
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
			},
		},
	}
}
