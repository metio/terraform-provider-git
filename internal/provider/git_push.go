/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ssh2 "golang.org/x/crypto/ssh"
)

func CreatePushOptions(ctx context.Context, inputs *PushResourceModel, diag *diag.Diagnostics) *git.PushOptions {
	options := &git.PushOptions{}

	if len(inputs.RefSpecs.Elements()) > 0 {
		refSpecs := make([]config.RefSpec, len(inputs.RefSpecs.Elements()))
		diag.Append(inputs.RefSpecs.ElementsAs(ctx, &refSpecs, false)...)
		if diag.HasError() {
			return nil
		}
		options.RefSpecs = refSpecs
		tflog.Trace(ctx, "using 'RefSpecs'", map[string]interface{}{
			"RefSpecs": refSpecs,
		})
	} else {
		return nil
	}

	options.RemoteName = inputs.Remote.ValueString()
	tflog.Trace(ctx, "using 'RemoteName'", map[string]interface{}{
		"RemoteName": inputs.Remote.ValueString(),
	})

	options.Prune = inputs.Prune.ValueBool()
	tflog.Trace(ctx, "using 'Prune'", map[string]interface{}{
		"Prune": inputs.Prune.ValueBool(),
	})

	options.Force = inputs.Force.ValueBool()
	tflog.Trace(ctx, "using 'Force'", map[string]interface{}{
		"Force": inputs.Force.ValueBool(),
	})

	if !inputs.Auth.IsNull() {
		basicAuth, basicOk := inputs.Auth.Attributes()["basic"].(types.Object)
		bearerAuth, bearerOk := inputs.Auth.Attributes()["bearer"].(types.String)
		sshKeyAuth, sshKeyOk := inputs.Auth.Attributes()["ssh_key"].(types.Object)
		sshAgentAuth, sshAgentOk := inputs.Auth.Attributes()["ssh_agent"].(types.Object)
		sshPasswordAuth, sshPasswordOk := inputs.Auth.Attributes()["ssh_password"].(types.Object)

		if basicOk && !basicAuth.IsNull() {
			username := basicAuth.Attributes()["username"].(types.String)
			password := basicAuth.Attributes()["password"].(types.String)

			options.Auth = &http.BasicAuth{
				Username: username.ValueString(),
				Password: password.ValueString(),
			}
		} else if bearerOk && !bearerAuth.IsNull() {
			options.Auth = &http.TokenAuth{
				Token: bearerAuth.ValueString(),
			}
		} else if sshKeyOk && !sshKeyAuth.IsNull() {
			username := sshKeyAuth.Attributes()["username"].(types.String)
			password := sshKeyAuth.Attributes()["password"].(types.String)

			var sshKeys *ssh.PublicKeys
			var err error
			if sshKeyAuth.Attributes()["private_key_path"] != nil {
				keyPath := sshKeyAuth.Attributes()["private_key_path"].(types.String)
				sshKeys, err = ssh.NewPublicKeysFromFile(username.ValueString(), keyPath.ValueString(), password.ValueString())
			} else if sshKeyAuth.Attributes()["private_key_pem"] != nil {
				keyPem := sshKeyAuth.Attributes()["private_key_pem"].(types.String)
				sshKeys, err = ssh.NewPublicKeys(username.ValueString(), []byte(keyPem.ValueString()), password.ValueString())
			} else {
				diag.AddError(
					"Invalid SSH key configuration",
					"Either path or PEM data must be specified",
				)
				return nil
			}
			if err != nil {
				diag.AddError(
					"Cannot use given SSH configuration",
					"SSH configuration failed because of: "+err.Error(),
				)
				return nil
			}

			if callback := knownHostsCallback(ctx, sshKeyAuth, diag); callback != nil {
				cfg, _ := sshKeys.ClientConfig()
				cfg.HostKeyCallback = callback
			}

			options.Auth = sshKeys
		} else if sshAgentOk && !sshAgentAuth.IsNull() {
			username := sshAgentAuth.Attributes()["username"].(types.String)

			agentAuth, err := ssh.NewSSHAgentAuth(username.ValueString())
			if err != nil {
				diag.AddError(
					"Cannot use SSH agent authentication",
					"Using SSH agent failed because of: "+err.Error(),
				)
				return nil
			}

			if callback := knownHostsCallback(ctx, sshAgentAuth, diag); callback != nil {
				cfg, _ := agentAuth.ClientConfig()
				cfg.HostKeyCallback = callback
			}

			options.Auth = agentAuth
		} else if sshPasswordOk && !sshPasswordAuth.IsNull() {
			username := sshPasswordAuth.Attributes()["username"].(types.String)
			password := sshPasswordAuth.Attributes()["password"].(types.String)

			passwordAuth := &ssh.Password{
				User:     username.ValueString(),
				Password: password.ValueString(),
			}

			if callback := knownHostsCallback(ctx, sshPasswordAuth, diag); callback != nil {
				cfg, _ := passwordAuth.ClientConfig()
				cfg.HostKeyCallback = callback
			}

			options.Auth = passwordAuth
		}
	}

	return options
}

func knownHostsCallback(ctx context.Context, object types.Object, diag *diag.Diagnostics) ssh2.HostKeyCallback {
	var files []string
	knownHosts, ok := object.Attributes()["known_hosts"].(types.Set)
	if ok && !knownHosts.IsNull() {
		diag.Append(knownHosts.ElementsAs(ctx, &files, false)...)
		if diag.HasError() {
			return nil
		}
	}
	callback, err := ssh.NewKnownHostsCallback(files...)
	if err != nil {
		diag.AddWarning(
			"Cannot use given known hosts - ",
			"Known hosts configuration failed because of: "+err.Error(),
		)
		return ssh2.InsecureIgnoreHostKey()
	}
	return callback
}
