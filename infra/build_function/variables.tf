variable "name" {
  description = "value of the name of the builder"
  type        = string
}
variable "region" {
  description = "value of the region of the builder"
  type        = string
}
variable "source_repo" {
  description = "value of the source_repo of the builder"
  type        = string
}
variable "function_path" {
  description = "value of the function_path of the builder"
  type        = string
}
variable "trigger_topic" {
  description = "value of the trigger_topic of the builder"
  type        = string
}
variable "entry_point" {
  description = "value of the entry_point of the builder"
  type        = string
}