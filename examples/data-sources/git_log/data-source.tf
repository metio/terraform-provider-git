# get commit logs of repository
data "git_log" "log" {
  directory = "/path/to/git/repository"
}

# filter commits to certain paths
data "git_log" "filtered" {
  directory    = "/path/to/git/repository"
  filter_paths = ["some/directory/in/repository/*", "README.md"]
}

# read commit info for each commit from log
data "git_commit" "from_log" {
  for_each  = toset(data.git_log.filtered.commits)
  directory = "/path/to/git/repository"
  revision  = each.value
}
