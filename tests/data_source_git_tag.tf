# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_tag" "every" {
  for_each  = data.git_tags.all.tags
  directory = data.git_repository.repository.directory
  name      = each.key
}

output "data_source_git_tag_every" {
  value = data.git_tag.every[*]
}
