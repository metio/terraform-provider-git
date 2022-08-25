/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func signatureToObject(signature *object.Signature) types.Object {
	data := make(map[string]attr.Value)

	if signature != nil {
		data["name"] = types.String{Value: signature.Name}
		data["email"] = types.String{Value: signature.Email}
		data["timestamp"] = types.String{Value: signature.When.String()}
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
