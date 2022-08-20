terraform {
  required_providers {
    git = {
      source = "localhost/metio/git"
      version = "9999.99.99"
    }
  }
}

provider "git" {
  # Configuration options
}
