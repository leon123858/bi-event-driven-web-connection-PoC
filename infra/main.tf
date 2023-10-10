provider "google" {
  project = var.project_id
  region  = var.region
}

// source code repository
module "repository" {
  source     = "./repository"
  name       = "todolist"
  region     = var.region
  git_url    = var.git_url
  git_app_id = var.git_app_id
  git_token  = var.git_token
}

// cloud run default server
module "http-server" {
  source       = "./run"
  service_name = "http"
  region       = var.region
}

module "websocket-server" {
  source       = "./run"
  service_name = "websocket"
  region       = var.region
}

module "web-server" {
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
  service_name = "get_todolist"
  region       = var.region
  bucket_name  = google_storage_bucket.function-source.name
}

module "update_todolist" {
  source       = "./function"
  service_name = "update_todolist"
  region       = var.region
  bucket_name  = google_storage_bucket.function-source.name
}

module "create_todolist" {
  source       = "./function"
  service_name = "create_todolist"
  region       = var.region
  bucket_name  = google_storage_bucket.function-source.name
}

module "delete_todolist" {
  source       = "./function"
  service_name = "delete_todolist"
  region       = var.region
  bucket_name  = google_storage_bucket.function-source.name
}

// cloud build set
data "google_project" "project" {
}

locals {
  cloud_build_sa =  "serviceAccount:${data.google_project.project.number}@cloudbuild.gserviceaccount.com"
}

resource "google_project_iam_member" "build_r1" {
  project = var.project_id
  role    = "roles/run.admin"
  # service accout in prject named "Cloud Build Service Account"
  member  = local.cloud_build_sa
}

resource "google_project_iam_member" "build_r2" {
  project    = var.project_id
  role       = "roles/iam.serviceAccountUser"
  # service accout in prject named "Cloud Build Service Account"
  member  = local.cloud_build_sa

  depends_on = [google_project_iam_member.build_r1]
}

// cloud build run

// cloud build function