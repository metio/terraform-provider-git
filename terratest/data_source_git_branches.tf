# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_branches" "all" {
  directory = data.git_repository.repository.directory
}

output "data_source_git_branches_branches" {
  value = data.git_branches.all
}
