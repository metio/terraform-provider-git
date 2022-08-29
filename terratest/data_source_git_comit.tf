# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_commit" "current_head" {
  directory = data.git_repository.repository.directory
  revision  = "HEAD"
}

data "git_commit" "from_log" {
  for_each  = toset(data.git_log.from_head.commits)
  directory = data.git_repository.repository.directory
  revision  = each.value
}

output "data_source_git_commit_current_head" {
  value = data.git_commit.current_head
}

output "data_source_git_commit_from_log" {
  value = data.git_commit.from_log
}
