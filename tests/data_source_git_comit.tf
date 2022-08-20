# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_commit" "main_head" {
  directory = data.git_repository.repository.directory
  sha1      = data.git_branch.main.sha1
}

output "data_source_git_commit_origin" {
  value = data.git_commit.main_head
}
