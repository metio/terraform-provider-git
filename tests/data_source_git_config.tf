# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_config" "defaults" {
  directory = data.git_repository.repository.directory
}

data "git_config" "local" {
  directory = data.git_repository.repository.directory
  scope     = "local"
}

data "git_config" "global" {
  directory = data.git_repository.repository.directory
  scope     = "global"
}

data "git_config" "system" {
  directory = data.git_repository.repository.directory
  scope     = "system"
}

output "data_source_git_config_defaults" {
  value = data.git_config.defaults
}

output "data_source_git_config_local" {
  value = data.git_config.local
}

output "data_source_git_config_global" {
  value = data.git_config.global
}

output "data_source_git_config_system" {
  value = data.git_config.system
}
