/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type GitTag struct {
	Annotated   types.Bool   `tfsdk:"annotated"`
	Lightweight types.Bool   `tfsdk:"lightweight"`
	SHA1        types.String `tfsdk:"sha1"`
}

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