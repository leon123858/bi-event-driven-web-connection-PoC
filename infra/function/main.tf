resource "google_pubsub_topic" "default" {
  name = "${var.service_name}-topic"
}

resource "google_storage_bucket_object" "object" {
  name   = "${var.service_name}-source.zip"
  bucket = var.bucket_name
  source = "${path.module}/default-source.zip" # Add path to the zipped function source code
}

resource "google_cloudfunctions2_function" "default" {
  name        = var.service_name
  location    = var.region
  description = "a default function"

  build_config {
    runtime     = "go119"
    entry_point = "HelloPubSub" # Set the entry point
    source {
      storage_source {
        bucket = var.bucket_name
        object = google_storage_bucket_object.object.name
      }
    }
  }

  service_config {
    max_instance_count             = 3
    min_instance_count             = 0
    timeout_seconds                = 60
    ingress_settings               = "ALLOW_INTERNAL_ONLY"
    all_traffic_on_latest_revision = true
  }

  event_trigger {
    trigger_region = var.region
    event_type     = "google.cloud.pubsub.topic.v1.messagePublished"
    pubsub_topic   = google_pubsub_topic.default.id
    retry_policy   = "RETRY_POLICY_DO_NOT_RETRY"
  }
}