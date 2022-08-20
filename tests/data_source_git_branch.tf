# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_branch" "main" {
  directory = data.git_repository.repository.directory
  name      = "main"
}

data "git_branch" "current" {
  directory = data.git_repository.repository.directory
  name      = data.git_repository.repository.branch
}

output "data_source_git_branch_main" {
  value = data.git_branch.main
}

output "data_source_git_branch_current" {
  value = data.git_branch.current
}
