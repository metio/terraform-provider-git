/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"
)

func signatureToObject(signature *object.Signature) types.Object {
	data := make(map[string]attr.Value)

	if signature != nil {
		data["name"] = types.StringValue(signature.Name)
		data["email"] = types.StringValue(signature.Email)
		data["timestamp"] = types.StringValue(signature.When.Format(time.RFC3339))
	} else {
		data["name"] = types.StringNull()
		data["email"] = types.StringNull()
		data["timestamp"] = types.StringNull()
	}

	return types.ObjectValueMust(
		map[string]attr.Type{
			"name":      types.StringType,
			"email":     types.StringType,
			"timestamp": types.StringType,
		},
		data,
	)
}

func signatureToObjectWithoutTimestamp(signature *object.Signature) types.Object {
	data := make(map[string]attr.Value)

	data["name"] = types.StringNull()
	data["email"] = types.StringNull()

	if signature != nil {
		if signature.Name != "" {
			data["name"] = types.StringValue(signature.Name)
		}
		if signature.Email != "" {
			data["email"] = types.StringValue(signature.Email)
		}
	}

	return types.ObjectValueMust(
		map[string]attr.Type{
			"name":  types.StringType,
			"email": types.StringType,
		},
		data,
	)
}

func createCommitOptions(ctx context.Context, inputs commitResourceModel) *git.CommitOptions {
	options := &git.CommitOptions{}

	options.All = inputs.All.ValueBool()
	tflog.Trace(ctx, "using 'All'", map[string]interface{}{
		"all": options.All,
	})

	options.AllowEmptyCommits = inputs.AllowEmptyCommits.ValueBool()
	tflog.Trace(ctx, "using 'AllowEmptyCommits'", map[string]interface{}{
		"allow empty commits": options.AllowEmptyCommits,
	})

	if !inputs.Author.IsNull() {
		options.Author = objectToSignature(&inputs.Author)
		tflog.Trace(ctx, "using 'Author'", map[string]interface{}{
			"name":  options.Author.Name,
			"email": options.Author.Email,
		})
	}

	if !inputs.Committer.IsNull() {
		options.Committer = objectToSignature(&inputs.Committer)
		tflog.Trace(ctx, "using 'Committer'", map[string]interface{}{
			"name":  options.Committer.Name,
			"email": options.Committer.Email,
		})
	}

	return options
}

func objectToSignature(obj *types.Object) *object.Signature {
	sig := &object.Signature{When: time.Now()}

	name := obj.Attributes()["name"].(types.String)
	if !name.IsNull() {
		sig.Name = name.ValueString()
	}

	email := obj.Attributes()["email"].(types.String)
	if !email.IsNull() {
		sig.Email = email.ValueString()
	}

	return sig
}

func extractModifiedFiles(commit *object.Commit) []string {
	fileNames := make([]string, 0)
	files, err := commit.Files()
	if err != nil {
		return fileNames
	}
	err = files.ForEach(func(file *object.File) error {
		fileNames = append(fileNames, file.Name)
		return nil
	})
	if err != nil {
		return fileNames
	}
	return fileNames
}
