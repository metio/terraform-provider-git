# commit changes and supply message
resource "git_commit" "commit" {
  directory = "/path/to/git/repository"
  message   = "committed with terraform"
}

# specify author
resource "git_commit" "author" {
  directory = "/path/to/git/repository"
  message   = "committed with terraform"
  author    = {
    name  = "terraform"
    email = "automation@example.com"
  }
}

# specify committer
resource "git_commit" "committer" {
  directory = "/path/to/git/repository"
  message   = "committed with terraform"
  committer = {
    name  = "terraform"
    email = "automation@example.com"
  }
}
