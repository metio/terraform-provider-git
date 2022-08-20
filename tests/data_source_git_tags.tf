data "git_tags" "tags" {
  directory = data.git_repository.repository.directory
}

output "data_source_git_tags_tags" {
  value = data.git_tags.tags
}
