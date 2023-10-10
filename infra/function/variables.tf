variable "service_name" {
  description = "The name of cloud run"
  type        = string
  default     = "docker-sample"
}
variable "region" {
  description = "The region of cloud run"
  type        = string
}
variable "bucket_name" {
  description = "value of bucket name save function source"
  type        = string
}
