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

data "git_repository" "repository" {
  directory = var.directory
}

data "git_branch" "current" {
  directory = data.git_repository.repository.directory
  name      = data.git_repository.repository.branch
}
