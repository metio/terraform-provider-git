# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_status" "status" {
  for_each  = data.git_statuses.statuses.files
  directory = data.git_repository.repository.directory
  file      = each.key
}

output "data_source_git_status_status" {
  value = data.git_status.status[*]
}
