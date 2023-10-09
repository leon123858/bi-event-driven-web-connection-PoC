variable "service_name" {
  description = "The name of cloud run"
  type        = string
  default     = "docker-sample"
}
variable "region" {
  description = "The region of cloud run"
  type        = string
}
variable "run-sa" {
  description = "The service account that cloud run used"
  type        = string
}
variable "db-name" {
  description = "The name of db"
  type        = string
}
