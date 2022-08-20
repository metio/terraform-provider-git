# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_remote" "origin" {
  directory = data.git_repository.repository.directory
  name      = "origin"
}

output "data_source_git_remote_origin" {
  value = data.git_remote.origin
}
