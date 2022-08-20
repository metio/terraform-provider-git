data "git_repository" "repository" {
  directory = var.git_repo_path
}

output "data_source_git_repository_repository" {
  value = data.git_repository.repository
}
