# push changes
resource "git_push" "push" {
  directory = "/path/to/git/repository"
  refspecs  = ["refs/heads/master:refs/heads/master"]
}

# specify remote
resource "git_push" "remote" {
  directory = "/path/to/git/repository"
  refspecs  = ["refs/heads/master:refs/heads/master"]
  remote    = "upstream"
}

# force push
resource "git_push" "remote" {
  directory = "/path/to/git/repository"
  refspecs  = ["refs/heads/master:refs/heads/master"]
  force     = true
}

# push on new commits
resource "git_commit" "commit" {
  directory = "/path/to/git/repository"
  message   = "committed with terraform"
}
resource "git_push" "push_on_commit" {
  directory = "/path/to/git/repository"
  refspecs  = ["refs/heads/master:refs/heads/master"]

  lifecycle {
    replace_triggered_by = [git_commit.commit.id]
  }
}
