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

# push with basic auth
resource "git_push" "remote" {
  directory = "/path/to/git/repository"
  refspecs  = ["refs/heads/master:refs/heads/master"]

  auth = {
    basic = {
      username = "some-username"
      password = "topsecret123"
    }
  }
}

# push with HTTP bearer token
resource "git_push" "remote" {
  directory = "/path/to/git/repository"
  refspecs  = ["refs/heads/master:refs/heads/master"]

  auth = {
    bearer = "some bearer token here"
  }
}

# push with SSH key
resource "git_push" "remote" {
  directory = "/path/to/git/repository"
  refspecs  = ["refs/heads/master:refs/heads/master"]

  auth = {
    ssh_key = {
      private_key_path = pathexpand("~/.ssh/id_rsa")
    }
  }
}

# push with SSH password
resource "git_push" "remote" {
  directory = "/path/to/git/repository"
  refspecs  = ["refs/heads/master:refs/heads/master"]

  auth = {
    ssh_password = {
      username = "bob"
      password = "hunter2"
    }
  }
}

# push with SSH agent
resource "git_push" "remote" {
  directory = "/path/to/git/repository"
  refspecs  = ["refs/heads/master:refs/heads/master"]

  auth = {
    ssh_agent = {}
  }
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
