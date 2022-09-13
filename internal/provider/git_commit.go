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
		data["name"] = types.String{Value: signature.Name}
		data["email"] = types.String{Value: signature.Email}
		data["timestamp"] = types.String{Value: signature.When.Format(time.RFC3339)}
	} else {
		data["name"] = types.String{Null: true}
		data["email"] = types.String{Null: true}
		data["timestamp"] = types.String{Null: true}
	}

	return types.Object{
		AttrTypes: map[string]attr.Type{
			"name":      types.StringType,
			"email":     types.StringType,
			"timestamp": types.StringType,
		},
		Attrs: data,
	}
}

func signatureToObjectWithoutTimestamp(signature *object.Signature) types.Object {
	data := make(map[string]attr.Value)

	data["name"] = types.String{Null: true}
	data["email"] = types.String{Null: true}

	if signature != nil {
		if signature.Name != "" {
			data["name"] = types.String{Value: signature.Name}
		}
		if signature.Email != "" {
			data["email"] = types.String{Value: signature.Email}
		}
	}

	return types.Object{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"email": types.StringType,
		},
		Attrs: data,
	}
}

func createCommitOptions(ctx context.Context, inputs commitResourceModel) *git.CommitOptions {
	options := &git.CommitOptions{}

	options.All = inputs.All.Value
	tflog.Trace(ctx, "using 'All'", map[string]interface{}{
		"all": options.All,
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

	name := obj.Attrs["name"].(types.String)
	if !name.IsNull() {
		sig.Name = name.Value
	}

	email := obj.Attrs["email"].(types.String)
	if !email.IsNull() {
		sig.Email = email.Value
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
