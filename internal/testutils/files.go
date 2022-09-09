//go:build testing

/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package testutils

import (
	"os"
	"path/filepath"
	"testing"
)

func TemporaryDirectory(t *testing.T) string {
	directory, err := os.MkdirTemp("", "terraform-provider-git")
	if err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(directory)
}

func WriteFile(t *testing.T, name string) {
	WriteFileContent(t, name, "hello world!")
}

func WriteFileContent(t *testing.T, name string, content string) {
	err := os.WriteFile(name, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
}
