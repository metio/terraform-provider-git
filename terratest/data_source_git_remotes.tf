# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_remotes" "remotes" {
  directory = data.git_repository.repository.directory
}

output "data_source_git_remotes_remotes" {
  value = data.git_remotes.remotes
}
