data "google_project" "project" {
}

// Create a secret containing the personal access token and grant permissions to the Service Agent
resource "google_secret_manager_secret" "github_token_secret" {
  secret_id = "${var.name}-secret"
  replication {
    user_managed {
      replicas {
        location = var.region
      }
    }
  }
}

resource "google_secret_manager_secret_version" "github_token_secret_version" {
  secret      = google_secret_manager_secret.github_token_secret.id
  secret_data = var.git_token
}

resource "google_secret_manager_secret_iam_member" "policy" {
  secret_id = google_secret_manager_secret.github_token_secret.id
  role      = "roles/secretmanager.secretAccessor"
  # account name is "Cloud Build Service Agent"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloudbuild.iam.gserviceaccount.com"

  depends_on = [google_secret_manager_secret_version.github_token_secret_version]
}

// Create the GitHub connection
resource "google_cloudbuildv2_connection" "my_connection" {
  location = var.region
  name     = "${var.name}-connection"

  github_config {
    app_installation_id = var.git_app_id
    authorizer_credential {
      oauth_token_secret_version = google_secret_manager_secret_version.github_token_secret_version.id
    }
  }

  depends_on = [google_secret_manager_secret_iam_member.policy]
}

resource "google_cloudbuildv2_repository" "my-repository" {
  name              = var.name
  location          = var.region
  parent_connection = google_cloudbuildv2_connection.my_connection.id
  remote_uri        = "${var.git_url}.git"
}