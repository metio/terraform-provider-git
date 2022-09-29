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
	"golang.org/x/crypto/ssh/knownhosts"
)

func CreatePushOptions(ctx context.Context, inputs *PushResourceModel, diag *diag.Diagnostics) *git.PushOptions {
	options := &git.PushOptions{}

	if len(inputs.RefSpecs.Elems) > 0 {
		refSpecs := make([]config.RefSpec, len(inputs.RefSpecs.Elems))
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

	options.RemoteName = inputs.Remote.Value
	tflog.Trace(ctx, "using 'RemoteName'", map[string]interface{}{
		"RemoteName": inputs.Remote.Value,
	})

	options.Prune = inputs.Prune.Value
	tflog.Trace(ctx, "using 'Prune'", map[string]interface{}{
		"Prune": inputs.Prune.Value,
	})

	options.Force = inputs.Force.Value
	tflog.Trace(ctx, "using 'Force'", map[string]interface{}{
		"Force": inputs.Force.Value,
	})

	if !inputs.Auth.IsNull() {
		basicAuth, basicOk := inputs.Auth.Attrs["basic"].(types.Object)
		bearerAuth, bearerOk := inputs.Auth.Attrs["bearer"].(types.String)
		sshKeyAuth, sshKeyOk := inputs.Auth.Attrs["ssh_key"].(types.Object)
		sshAgentAuth, sshAgentOk := inputs.Auth.Attrs["ssh_agent"].(types.Object)
		sshPasswordAuth, sshPasswordOk := inputs.Auth.Attrs["ssh_password"].(types.Object)

		if basicOk && !basicAuth.IsNull() {
			username := basicAuth.Attrs["username"].(types.String)
			password := basicAuth.Attrs["password"].(types.String)

			options.Auth = &http.BasicAuth{
				Username: username.Value,
				Password: password.Value,
			}
		} else if bearerOk && !bearerAuth.IsNull() {
			options.Auth = &http.TokenAuth{
				Token: bearerAuth.Value,
			}
		} else if sshKeyOk && !sshKeyAuth.IsNull() {
			username := sshKeyAuth.Attrs["username"].(types.String)
			password := sshKeyAuth.Attrs["password"].(types.String)

			var sshKeys *ssh.PublicKeys
			var err error
			if sshKeyAuth.Attrs["private_key_path"] != nil {
				keyPath := sshKeyAuth.Attrs["private_key_path"].(types.String)
				sshKeys, err = ssh.NewPublicKeysFromFile(username.Value, keyPath.Value, password.Value)
			} else if sshKeyAuth.Attrs["private_key_pem"] != nil {
				keyPem := sshKeyAuth.Attrs["private_key_pem"].(types.String)
				sshKeys, err = ssh.NewPublicKeys(username.Value, []byte(keyPem.Value), password.Value)
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
				sshKeys.HostKeyCallback = callback
			}

			options.Auth = sshKeys
		} else if sshAgentOk && !sshAgentAuth.IsNull() {
			username := sshAgentAuth.Attrs["username"].(types.String)

			agentAuth, err := ssh.NewSSHAgentAuth(username.Value)
			if err != nil {
				diag.AddError(
					"Cannot use SSH agent authentication",
					"Using SSH agent failed because of: "+err.Error(),
				)
				return nil
			}

			if callback := knownHostsCallback(ctx, sshAgentAuth, diag); callback != nil {
				agentAuth.HostKeyCallback = callback
			}

			options.Auth = agentAuth
		} else if sshPasswordOk && !sshPasswordAuth.IsNull() {
			username := sshPasswordAuth.Attrs["username"].(types.String)
			password := sshPasswordAuth.Attrs["password"].(types.String)

			passwordAuth := &ssh.Password{
				User:     username.Value,
				Password: password.Value,
			}

			if callback := knownHostsCallback(ctx, sshPasswordAuth, diag); callback != nil {
				passwordAuth.HostKeyCallback = callback
			}

			options.Auth = passwordAuth
		}
	}

	return options
}

func knownHostsCallback(ctx context.Context, object types.Object, diag *diag.Diagnostics) ssh2.HostKeyCallback {
	knownHosts, ok := object.Attrs["known_hosts"].(types.Set)
	if ok && !knownHosts.IsNull() {
		hosts := make([]string, len(knownHosts.Elems))
		diag.Append(knownHosts.ElementsAs(ctx, &hosts, false)...)
		if diag.HasError() {
			return nil
		}
		callback, err := knownhosts.New(hosts...)
		if err != nil {
			diag.AddError(
				"Cannot use given known hosts",
				"Known hosts configuration failed because of: "+err.Error(),
			)
			return nil
		}
		return callback
	}
	return nil
}
