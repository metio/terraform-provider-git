/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func createTagReference(repository *git.Repository, inputs resourceGitTagSchema) (*plumbing.Reference, error) {
	if inputs.SHA1.IsNull() || inputs.SHA1.IsUnknown() {
		head, err := repository.Head()
		if err != nil {
			return nil, err
		}
		return head, nil
	}

	return plumbing.NewHashReference("tag", plumbing.NewHash(inputs.SHA1.Value)), nil
}

func createOptions(inputs resourceGitTagSchema) *git.CreateTagOptions {
	if inputs.Message.IsNull() || inputs.Message.IsUnknown() {
		return nil
	}
	return &git.CreateTagOptions{
		Message: inputs.Message.Value,
	}
}
