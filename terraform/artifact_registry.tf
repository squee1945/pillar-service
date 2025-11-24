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

# This repository holds images that are used in the runner service.
resource "google_artifact_registry_repository" "runner_images" {
  project       = var.project_id
  location      = var.region
  repository_id = "runner-images"
  description   = "A repo to store all the runner images."
  format        = "DOCKER"
  provider      = google-beta
}

# This repository holds Go module artifacts that are built when the
# runner kicks off a "sub-build". These artifacts are generated and stored
# so that Artifact Analysis can generate provenance.
resource "google_artifact_registry_repository" "sub_build_go_repository" {
  project       = var.project_id
  location      = var.region
  repository_id = "sub-build-go-repository"
  format        = "GO"
  description   = "Repository for storing Go modules for the purposes of attestation generation."
  cleanup_policies {
    id     = "delete-old-go-modules"
    action = "DELETE"
    condition {
      older_than = "720h" # 30 days
    }
  }
}
