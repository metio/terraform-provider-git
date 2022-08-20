# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_branch" "current" {
  count     = data.git_repository.repository.branch == null ? 0 : 1
  directory = data.git_repository.repository.directory
  name      = data.git_repository.repository.branch
}

data "git_branch" "every" {
  for_each  = data.git_branches.all.branches
  directory = data.git_repository.repository.directory
  name      = each.key
}

output "data_source_git_branch_current" {
  value = data.git_branch.current[*]
}

output "data_source_git_branch_every" {
  value = data.git_branch.every[*]
}
