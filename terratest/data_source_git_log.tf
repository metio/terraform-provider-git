# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "git_log" "log" {
  directory = data.git_repository.repository.directory
  max_count = 5
}

data "git_log" "single_file" {
  directory    = data.git_repository.repository.directory
  filter_paths = ["tests/.gitignore"]
}

data "git_log" "multiple_files" {
  directory    = data.git_repository.repository.directory
  filter_paths = ["tests/.gitignore", "tests/data_source_git_branch.tf"]
}

data "git_log" "single_directory" {
  directory    = data.git_repository.repository.directory
  filter_paths = ["tests/*"]
}

data "git_log" "multiple_directories" {
  directory    = data.git_repository.repository.directory
  filter_paths = ["tests/*", "LICENSES/*"]
}

data "git_log" "all" {
  directory = data.git_repository.repository.directory
  all       = true
  max_count = 5
}

data "git_log" "from_head" {
  directory = data.git_repository.repository.directory
  from      = "HEAD"
  max_count = 5
}

data "git_log" "from_at" {
  directory = data.git_repository.repository.directory
  from      = "@"
  max_count = 5
}

data "git_log" "from_tag" {
  directory = data.git_repository.repository.directory
  from      = "2022.8.12"
  max_count = 5
}

data "git_log" "range" {
  directory = data.git_repository.repository.directory
  since     = timeadd(timestamp(), "-168h")
  until     = timeadd(timestamp(), "-24h")
}

data "git_log" "skipped" {
  directory = data.git_repository.repository.directory
  max_count = 5
  skip      = 3
}

data "git_log" "time_ordered" {
  directory = data.git_repository.repository.directory
  since     = timeadd(timestamp(), "-168h")
  until     = timeadd(timestamp(), "-24h")
  order     = "time"
}

data "git_log" "depth_ordered" {
  directory = data.git_repository.repository.directory
  since     = timeadd(timestamp(), "-168h")
  until     = timeadd(timestamp(), "-24h")
  order     = "depth"
}

data "git_log" "breadth_ordered" {
  directory = data.git_repository.repository.directory
  since     = timeadd(timestamp(), "-168h")
  until     = timeadd(timestamp(), "-24h")
  order     = "breadth"
}

output "data_source_git_log_log" {
  value = data.git_log.log
}

output "data_source_git_log_single_file" {
  value = data.git_log.single_file
}

output "data_source_git_log_multiple_files" {
  value = data.git_log.multiple_files
}

output "data_source_git_log_single_directory" {
  value = data.git_log.single_directory
}

output "data_source_git_log_multiple_directories" {
  value = data.git_log.multiple_directories
}

output "data_source_git_log_all" {
  value = data.git_log.all
}

output "data_source_git_log_from_head" {
  value = data.git_log.from_head
}

output "data_source_git_log_from_at" {
  value = data.git_log.from_at
}

output "data_source_git_log_from_tag" {
  value = data.git_log.from_tag
}

output "data_source_git_log_range" {
  value = data.git_log.range
}

output "data_source_git_log_skipped" {
  value = data.git_log.skipped
}

output "data_source_git_log_time_ordered" {
  value = data.git_log.time_ordered
}

output "data_source_git_log_depth_ordered" {
  value = data.git_log.depth_ordered
}

output "data_source_git_log_breadth_ordered" {
  value = data.git_log.breadth_ordered
}
