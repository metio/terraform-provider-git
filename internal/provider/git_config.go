/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import "github.com/go-git/go-git/v5/config"

func mapConfigScope(userInput string) config.Scope {
	switch userInput {
	case "local":
		return config.LocalScope
	case "system":
		return config.SystemScope
	default:
		return config.GlobalScope
	}
}
