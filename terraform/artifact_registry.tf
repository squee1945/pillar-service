# This respoitory holds images that are for the Pillar service itself
# (i.e., the service that responds to webhook events).
resource "google_artifact_registry_repository" "pillar_service" {
  project       = var.project_id
  location      = var.region
  repository_id = "pillar-service"
  description   = "A repo to store all the Pillar Service images."
  format        = "DOCKER"
  provider      = google-beta
  depends_on    = [google_project_service.default["artifactregistry.googleapis.com"]]
}
