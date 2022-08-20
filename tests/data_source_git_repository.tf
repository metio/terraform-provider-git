# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_repository" "repository" {
  directory = var.git_repo_path
}

output "data_source_git_repository_repository" {
  value = data.git_repository.repository
}
