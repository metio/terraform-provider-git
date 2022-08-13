/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func getWorktree(repository *git.Repository, diag *diag.Diagnostics) (*git.Worktree, error) {
	worktree, err := repository.Worktree()
	if err == git.ErrIsBareRepository {
		return nil, nil
	} else if err == nil {
		return worktree, nil
	}
	diag.AddError(
		"Cannot read worktree",
		"Could not read worktree because of: "+err.Error(),
	)
	return nil, err
}
