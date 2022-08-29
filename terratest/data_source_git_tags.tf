# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_tags" "all" {
  directory = data.git_repository.repository.directory
}

output "data_source_git_tags_all" {
  value = data.git_tags.all
}
