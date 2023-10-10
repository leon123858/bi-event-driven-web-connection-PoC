resource "google_artifact_registry_repository" "docker-images" {
  location      = var.region
  repository_id = var.name
  format        = "DOCKER"

  docker_config {
    immutable_tags = false
  }
}

resource "google_cloudbuild_trigger" "manual-trigger" {
  location = var.region
  name     = "${var.name}-trigger"

  source_to_build {
    repository = var.source_repo
    ref        = "refs/heads/main"
    repo_type  = "GITHUB"
  }

  build {
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "--no-cache", "-t", "$_AR_HOSTNAME/$PROJECT_ID/$_REPO_NAME/$_SERVICE_NAME:$COMMIT_SHA", "${var.name}", "-f", "${var.docker_file_path}"]
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["push", "$_AR_HOSTNAME/$PROJECT_ID/$_REPO_NAME/$_SERVICE_NAME:$COMMIT_SHA"]
    }
    step {
      name = "gcr.io/cloud-builders/gcloud"
      args = ["run", "deploy", "${var.name}", "--image", "$_AR_HOSTNAME/$PROJECT_ID/$_REPO_NAME/$_SERVICE_NAME:$COMMIT_SHA", "--region", "${var.region}", "--platform", "managed", "--allow-unauthenticated"]
    }

    substitutions = {
      _SERVICE_NAME : var.name
      _DEPLOY_REGION : var.region
      _REPO_NAME : google_artifact_registry_repository.docker-images.name
      _AR_HOSTNAME : "${var.region}-docker.pkg.dev"
    }
  }


  // If this is set on a build, it will become pending when it is run, 
  // and will need to be explicitly approved to start.
  approval_config {
    approval_required = false
  }
}