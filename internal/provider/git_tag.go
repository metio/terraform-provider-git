/*
 * SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
 * SPDX-License-Identifier: BSD0
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

func createTagOptions(inputs resourceGitTagSchema) *git.CreateTagOptions {
	if inputs.Message.IsNull() || inputs.Message.IsUnknown() {
		return nil
	}
	return &git.CreateTagOptions{
		Message: inputs.Message.Value,
	}
}

func getTagReference(ctx context.Context, repository *git.Repository, tagName string, diag *diag.Diagnostics) *plumbing.Reference {
	tag, err := repository.Tag(tagName)
	if err != nil {
		diag.AddError(
			"Cannot read tag",
			"Could not read tag ["+tagName+"] because of: "+err.Error(),
		)
		return nil
	}
	tflog.Trace(ctx, "read tag", map[string]interface{}{
		"tag": tag,
	})
	return tag
}

func getTagObject(ctx context.Context, repository *git.Repository, hash plumbing.Hash, diag *diag.Diagnostics) (*object.Tag, error) {
	tag, err := repository.TagObject(hash)
	if err == plumbing.ErrObjectNotFound {
		tflog.Trace(ctx, "lightweight tag", map[string]interface{}{
			"hash": hash.String(),
		})
		return nil, nil
	} else if err == nil {
		tflog.Trace(ctx, "annotated tag", map[string]interface{}{
			"hash": hash.String(),
		})
		return tag, nil
	} else {
		diag.AddError(
			"Cannot read tag",
			"Could not read tag at hash ["+hash.String()+"] because of: "+err.Error(),
		)
		return nil, err
	}
}
