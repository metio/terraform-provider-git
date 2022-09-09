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

resource "git_add" "add" {
  directory = var.directory
  add_paths = var.add_paths
}

data "git_statuses" "multiple_files" {
  directory = git_add.add.directory
}
