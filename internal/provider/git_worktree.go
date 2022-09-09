/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func createCommit(worktree *git.Worktree, message string, options *git.CommitOptions, diag *diag.Diagnostics) *plumbing.Hash {
	hash, err := worktree.Commit(message, options)
	if err != nil {
		diag.AddError(
			"Cannot create commit",
			"Could not create commit because of: "+err.Error(),
		)
		return nil
	}
	return &hash
}

func getStatus(ctx context.Context, worktree *git.Worktree, diag *diag.Diagnostics) git.Status {
	status, err := worktree.Status()
	if err != nil {
		diag.AddError(
			"Cannot read status",
			"Could not read status because of: "+err.Error(),
		)
		return nil
	}
	tflog.Trace(ctx, "read status", map[string]interface{}{
		"status": status.String(),
	})
	return status
}
