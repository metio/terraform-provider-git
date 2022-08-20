data "git_remotes" "remotes" {
  directory = data.git_repository.repository.directory
}

output "data_source_git_remotes_remotes" {
  value = data.git_remotes.remotes
}
