/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func resolveRevision(ctx context.Context, repository *git.Repository, revision string, diag *diag.Diagnostics) *plumbing.Hash {
	hash, err := repository.ResolveRevision(plumbing.Revision(revision))
	if err != nil {
		diag.AddError(
			"Cannot resolve revision",
			"Could not resolve revision ["+revision+"] because of: "+err.Error(),
		)
		return nil
	}
	tflog.Trace(ctx, "resolved revision", map[string]interface{}{
		"revision": revision,
		"hash":     hash.String(),
	})
	return hash
}

func getCommit(ctx context.Context, repository *git.Repository, hash *plumbing.Hash, diag *diag.Diagnostics) *object.Commit {
	commitObject, err := repository.CommitObject(*hash)
	if err != nil {
		diag.AddError(
			"Cannot read commit",
			"Could not read commit ["+hash.String()+"] because of: "+err.Error(),
		)
		return nil
	}
	tflog.Trace(ctx, "read commit", map[string]interface{}{
		"hash": hash,
	})
	return commitObject
}
