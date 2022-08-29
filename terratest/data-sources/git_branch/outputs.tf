output "directory" {
  value = data.git_branch.branch.directory
}

output "name" {
  value = data.git_branch.branch.name
}

output "id" {
  value = data.git_branch.branch.id
}

output "rebase" {
  value = data.git_branch.branch.rebase
}

output "remote" {
  value = data.git_branch.branch.remote
}

output "sha1" {
  value = data.git_branch.branch.sha1
}
