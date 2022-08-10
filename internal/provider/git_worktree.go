/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
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
