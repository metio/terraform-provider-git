/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"io"
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

func addFile(worktree *git.Worktree, name string, diag *diag.Diagnostics) error {
	_, err := worktree.Add(name)
	if err != nil {
		diag.AddError(
			"Cannot add file to worktree",
			"The given file ["+name+"] cannot be added to the worktree because of: "+err.Error(),
		)
		return err
	}
	return nil
}

func readFileSha1(err error, worktree *git.Worktree, name string, diag *diag.Diagnostics) string {
	fileInfo, err := worktree.Filesystem.Lstat(name)
	if err != nil {
		diag.AddError(
			"Cannot get file infos",
			"Could not get infos about ["+name+"] because of: "+err.Error(),
		)
		return ""
	}
	if fileInfo.IsDir() {
		diag.AddError(
			"Cannot open directory",
			"The given path ["+name+"] is a directory but must be a file.",
		)
		return ""
	}
	file, err := worktree.Filesystem.Open(name)
	if err != nil {
		diag.AddError(
			"Cannot open file",
			"Could not open file ["+name+"] because of: "+err.Error(),
		)
		return ""
	}
	defer func(file billy.File) {
		if err := file.Close(); err != nil {
			diag.AddError(
				"Cannot close file handle",
				"Could not close file ["+name+"] because of: "+err.Error(),
			)
		}
	}(file)
	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		diag.AddError(
			"Cannot generate SHA1 hash",
			"Could not calculate SHA1 hash for file ["+name+"] because of: "+err.Error(),
		)
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}
