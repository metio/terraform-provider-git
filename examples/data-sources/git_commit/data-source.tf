# get commit by full sha1
data "git_commit" "full_sha1" {
  directory = "/path/to/git/repository"
  revision  = "dae86e1950b1277e545cee180551750029cfe735"
}

# get commit by short sha1
data "git_commit" "short_sha1" {
  directory = "/path/to/git/repository"
  revision  = "dae86e"
}

# get commit by refname
data "git_commit" "refname" {
  directory = "/path/to/git/repository"
  revision  = "main"
}

# get commit by head shortcut
data "git_commit" "head_shortcut" {
  directory = "/path/to/git/repository"
  revision  = "@"
}

# get commit by parent
data "git_commit" "head_parent" {
  directory = "/path/to/git/repository"
  revision  = "HEAD~1"
}
