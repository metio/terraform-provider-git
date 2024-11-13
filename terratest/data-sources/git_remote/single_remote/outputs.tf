# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_remote.single_remote.directory
}

output "id" {
  value = data.git_remote.single_remote.id
}

output "name" {
  value = data.git_remote.single_remote.name
}

output "urls" {
  value = data.git_remote.single_remote.urls
}

locals {
  selected_remote_url = data.git_remote.single_remote.urls[0]
  splitted_url = split("/", local.selected_remote_url)
  last_path_segment = local.splitted_url[length(local.splitted_url) - 1]
  basename = trimsuffix(local.last_path_segment, ".git")
}

output "upstream_repository_name" {
  value = local.basename
}
