/*
 * This file is part of terraform-gitProvider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-gitProvider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func extractGitRemoteUrls(remote *git.Remote) []types.String {
	var remoteUrls []types.String
	for _, url := range remote.Config().URLs {
		remoteUrls = append(remoteUrls, types.String{Value: url})
	}
	return remoteUrls
}
