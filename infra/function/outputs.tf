output "name" {
  value = google_cloudfunctions2_function.default.name
}
output "topic" {
  value = google_pubsub_topic.default.name
}