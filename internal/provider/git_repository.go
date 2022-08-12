/*
 * SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
 * SPDX-License-Identifier: BSD0
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func openRepository(ctx context.Context, directory string, diag *diag.Diagnostics) *git.Repository {
	repository, err := git.PlainOpen(directory)
	if err != nil {
		diag.AddError(
			"Cannot open repository",
			"Could not open git repository ["+directory+"] because of: "+err.Error(),
		)
		return nil
	}
	tflog.Trace(ctx, "opened repository", map[string]interface{}{
		"directory": directory,
	})
	return repository
}
