data "git_remote" "remote" {
  directory = "/path/to/git/repository"
  name      = "origin"
}

# ask the remote for its default branch, similar to 'git remote show origin'
data "git_remote" "default_branch" {
  directory    = "/path/to/git/repository"
  name         = "origin"
  query_remote = true
}

# query a private remote over SSH
data "git_remote" "private" {
  directory    = "/path/to/git/repository"
  name         = "origin"
  query_remote = true

  auth = {
    ssh_agent = {}
  }
}
