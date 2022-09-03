# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

terraform {
  required_providers {
    git = {
      source  = "localhost/metio/git"
      version = "9999.99.99"
    }
  }
}

provider "git" {
  # Configuration options
}

data "git_log" "from_head" {
  directory = var.directory
  from      = "HEAD"
  max_count = 5
}

data "git_commit" "from_log" {
  for_each  = toset(data.git_log.from_head.commits)
  directory = var.directory
  revision  = each.value
}
