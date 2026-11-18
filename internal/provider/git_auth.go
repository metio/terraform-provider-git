/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"net"
	"strconv"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ssh2 "golang.org/x/crypto/ssh"
)

func authOptions(ctx context.Context, auth types.Object, url string, diag *diag.Diagnostics) transport.AuthMethod {
	basicAuth, basicOk := auth.Attributes()["basic"].(types.Object)
	bearerAuth, bearerOk := auth.Attributes()["bearer"].(types.String)
	sshKeyAuth, sshKeyOk := auth.Attributes()["ssh_key"].(types.Object)
	sshAgentAuth, sshAgentOk := auth.Attributes()["ssh_agent"].(types.Object)
	sshPasswordAuth, sshPasswordOk := auth.Attributes()["ssh_password"].(types.Object)

	if basicOk && !basicAuth.IsNull() {
		username := basicAuth.Attributes()["username"].(types.String)
		password := basicAuth.Attributes()["password"].(types.String)

		return &http.BasicAuth{
			Username: username.ValueString(),
			Password: password.ValueString(),
		}
	} else if bearerOk && !bearerAuth.IsNull() {
		return &http.TokenAuth{
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

		if callback, algorithms := knownHostsCallback(ctx, sshKeyAuth, url, diag); callback != nil {
			sshKeys.HostKeyCallback = callback
			sshKeys.HostKeyAlgorithms = algorithms
		}
		if diag.HasError() {
			return nil
		}

		return sshKeys
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

		if callback, algorithms := knownHostsCallback(ctx, sshAgentAuth, url, diag); callback != nil {
			agentAuth.HostKeyCallback = callback
			agentAuth.HostKeyAlgorithms = algorithms
		}
		if diag.HasError() {
			return nil
		}

		return agentAuth
	} else if sshPasswordOk && !sshPasswordAuth.IsNull() {
		username := sshPasswordAuth.Attributes()["username"].(types.String)
		password := sshPasswordAuth.Attributes()["password"].(types.String)

		passwordAuth := &ssh.Password{
			User:     username.ValueString(),
			Password: password.ValueString(),
		}

		if callback, algorithms := knownHostsCallback(ctx, sshPasswordAuth, url, diag); callback != nil {
			passwordAuth.HostKeyCallback = callback
			passwordAuth.HostKeyAlgorithms = algorithms
		}
		if diag.HasError() {
			return nil
		}

		return passwordAuth
	}
	return nil
}

// knownHostsCallback verifies host keys against the given known_hosts files. Without any files, it returns nil so that
// go-git uses the default known_hosts files. A custom callback disables go-git's own detection of the host key
// algorithms, so they are looked up here for the host the remote URL points to.
func knownHostsCallback(ctx context.Context, object types.Object, url string, diag *diag.Diagnostics) (ssh2.HostKeyCallback, []string) {
	var files []string
	knownHosts, ok := object.Attributes()["known_hosts"].(types.Set)
	if ok && !knownHosts.IsNull() {
		diag.Append(knownHosts.ElementsAs(ctx, &files, false)...)
		if diag.HasError() {
			return nil, nil
		}
	}
	if len(files) == 0 {
		return nil, nil
	}

	db, err := ssh.NewKnownHostsDb(files...)
	if err != nil {
		diag.AddError(
			"Cannot use given known hosts",
			"Known hosts configuration failed because of: "+err.Error(),
		)
		return nil, nil
	}

	var algorithms []string
	if hostWithPort := sshHostWithPort(url); hostWithPort != "" {
		algorithms = db.HostKeyAlgorithms(hostWithPort)
	}
	return db.HostKeyCallback(), algorithms
}

// sshHostWithPort returns the address go-git dials for the given URL, honoring HostName and Port from the SSH config.
func sshHostWithPort(url string) string {
	endpoint, err := transport.NewEndpoint(url)
	if err != nil || endpoint.Host == "" {
		return ""
	}
	host := endpoint.Host
	port := endpoint.Port
	if ssh.DefaultSSHConfig != nil {
		if configHost := ssh.DefaultSSHConfig.Get(endpoint.Host, "Hostname"); configHost != "" {
			host = configHost
			if configPort, err := strconv.Atoi(ssh.DefaultSSHConfig.Get(endpoint.Host, "Port")); err == nil {
				port = configPort
			}
		}
	}
	if port <= 0 {
		port = ssh.DefaultPort
	}
	return net.JoinHostPort(host, strconv.Itoa(port))
}
