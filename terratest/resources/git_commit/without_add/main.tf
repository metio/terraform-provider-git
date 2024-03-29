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

resource "git_commit" "commit" {
  directory = var.directory
  message   = "committed with terraform"

  author = {
    name  = "terraform"
    email = "automation@example.com"
  }
}

data "git_commit" "commit" {
  directory = git_commit.commit.directory
  revision  = git_commit.commit.sha1
}

data "git_status" "single_file" {
  directory = git_commit.commit.directory
  file      = var.file
}
