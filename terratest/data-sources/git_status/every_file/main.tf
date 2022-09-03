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

data "git_statuses" "statuses" {
  directory = var.directory
}

data "git_status" "every_file" {
  for_each  = data.git_statuses.statuses.files
  directory = var.directory
  file      = each.key
}
