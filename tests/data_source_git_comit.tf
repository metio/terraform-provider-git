# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_commit" "current_head" {
  directory = data.git_repository.repository.directory
  revision  = "HEAD"
}

output "data_source_git_commit_current_head" {
  value = data.git_commit.current_head
}
