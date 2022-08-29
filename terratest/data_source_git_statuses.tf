# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_statuses" "statuses" {
  directory = data.git_repository.repository.directory
}

output "data_source_git_statuses_statuses" {
  value = data.git_statuses.statuses
}
