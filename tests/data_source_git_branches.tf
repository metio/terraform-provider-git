data "git_branches" "branches" {
  directory = data.git_repository.repository.directory
}

output "data_source_git_branches_branches" {
  value = data.git_branches.branches
}
