/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func mapConfigScope(userInput string) config.Scope {
	switch userInput {
	case "local":
		return config.LocalScope
	case "system":
		return config.SystemScope
	default:
		return config.GlobalScope
	}
}

func readConfig(ctx context.Context, repository *git.Repository, scope string, diag *diag.Diagnostics) *config.Config {
	cfg, err := repository.ConfigScoped(mapConfigScope(scope))
	if err != nil {
		diag.AddError(
			"Error reading config",
			"Could not read git config because of: "+err.Error(),
		)
		return nil
	}

	tflog.Trace(ctx, "read config", map[string]interface{}{
		"scope": scope,
	})
	return cfg
}
