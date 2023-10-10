provider "google" {
  project = var.project_id
  region  = var.region
}

// source code repository
module "repository" {
  source     = "./repository"
  name       = "todo-list"
  region     = var.region
  git_url    = var.git_url
  git_app_id = var.git_app_id
  git_token  = var.git_token
}

// cloud run default server
module "http_server" {
  source       = "./run"
  service_name = "http"
  region       = var.region
}

module "websocket_server" {
  source       = "./run"
  service_name = "websocket"
  region       = var.region
}

module "web_server" {
  source       = "./run"
  service_name = "web"
  region       = var.region
}

// cloud function default server
resource "google_storage_bucket" "function-source" {
  name                        = "${var.region}-${var.project_id}-gcf-source" # Every bucket name must be globally unique
  location                    = var.region
  uniform_bucket_level_access = true
}

module "get_todolist" {
  source       = "./function"
  service_name = "get-todolist"
  region       = var.region
  bucket_name  = google_storage_bucket.function-source.name
}

module "update_todolist" {
  source       = "./function"
  service_name = "update-todolist"
  region       = var.region
  bucket_name  = google_storage_bucket.function-source.name
}

module "create_todolist" {
  source       = "./function"
  service_name = "create-todolist"
  region       = var.region
  bucket_name  = google_storage_bucket.function-source.name
}

module "delete_todolist" {
  source       = "./function"
  service_name = "delete-todolist"
  region       = var.region
  bucket_name  = google_storage_bucket.function-source.name
}

// cloud build set
data "google_project" "project" {
}

locals {
  cloud_build_sa = "serviceAccount:${data.google_project.project.number}@cloudbuild.gserviceaccount.com"
}



resource "google_project_iam_member" "build_r0" {
  project = var.project_id
  role    = "roles/run.admin"
  # service accout in prject named "Cloud Build Service Account"
  member = local.cloud_build_sa
}

resource "google_project_iam_member" "build_r1" {
  project = var.project_id
  role    = "roles/cloudfunctions.developer"
  # service accout in prject named "Cloud Build Service Account"
  member = local.cloud_build_sa

  depends_on = [google_project_iam_member.build_r0]
}

resource "google_project_iam_member" "build_r2" {
  project = var.project_id
  role    = "roles/iam.serviceAccountUser"
  # service accout in prject named "Cloud Build Service Account"
  member = local.cloud_build_sa

  depends_on = [google_project_iam_member.build_r1]
}

// cloud build run
module "build_http" {
  source           = "./build_run"
  name             = module.http_server.name
  region           = var.region
  docker_file_path = "http/Dockerfile"
  source_repo      = module.repository.id
}

module "build_websocket" {
  source           = "./build_run"
  name             = module.websocket_server.name
  region           = var.region
  docker_file_path = "websocket/Dockerfile"
  source_repo      = module.repository.id
}

module "build_web" {
  source           = "./build_run"
  name             = module.web_server.name
  region           = var.region
  docker_file_path = "web/Dockerfile"
  source_repo      = module.repository.id
}

// cloud build function
module "build_get_todolist" {
  source       = "./build_function"
  name         = module.get_todolist.name
  region       = var.region
  function_path = "./functions"
  trigger_topic = module.get_todolist.topic
  entry_point   = "GetTodoList"
  source_repo   = module.repository.id
}
module "build_create_todolist" {
  source = "./build_function"
  name         = module.create_todolist.name
  region       = var.region
  function_path = "./functions"
  trigger_topic = module.create_todolist.topic
  entry_point   = "AddTodoItem"
  source_repo   = module.repository.id
}
module "build_update_todolist" {
  source = "./build_function"
  name         = module.update_todolist.name
  region       = var.region
  function_path = "./functions"
  trigger_topic = module.update_todolist.topic
  entry_point   = "UpdateTodoItem"
  source_repo   = module.repository.id
}
module "build_delete_todolist" {
  source = "./build_function"
  name         = module.delete_todolist.name
  region       = var.region
  function_path = "./functions"
  trigger_topic = module.delete_todolist.topic
  entry_point   = "RemoveTodoItem"
  source_repo   = module.repository.id
}
