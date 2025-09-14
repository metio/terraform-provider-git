//go:build testing

/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package testutils

import (
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func Signature() *object.Signature {
	return &object.Signature{
		Name:  "Some Person",
		Email: "person@example.com",
		When:  time.Now(),
	}
}
