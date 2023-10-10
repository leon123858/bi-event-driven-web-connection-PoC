variable "name" {
  description = "The name for repo"
  type        = string
}
variable "region" {
  default = "place replica secret"
  type    = string
}
variable "git_token" {
  description = "The key to access github repo. should have  repo and read:user permission, when code in org should have read:org permission optionally"
  type        = string
}
variable "git_app_id" {
  description = "id of the github app in repo. can find here https://github.com/settings/installations/<id>"
  type        = string
}
variable "git_url" {
  description = "url of the github repo"
  type        = string
}
