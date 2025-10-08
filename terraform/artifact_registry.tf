# resource "google_artifact_registry_repository" "runner_images" {
#   project       = var.project_id
#   location      = var.region
#   repository_id = "runner-images"
#   description   = "A repo to store all the runner images."
#   format        = "DOCKER"
#   provider      = google-beta
# }

resource "google_artifact_registry_repository" "pillar_service" {
  project       = var.project_id
  location      = var.region
  repository_id = "pillar-service"
  description   = "A repo to store all the Pillar Service images."
  format        = "DOCKER"
  provider      = google-beta
  depends_on    = [google_project_service.default["artifactregistry.googleapis.com"]]
}

# resource "google_artifact_registry_repository_iam_member" "sa_reader" {
#   project  = var.project_id
#   provider = google-beta

#   repository = google_artifact_registry_repository.runner_images.id
#   role       = "roles/artifactregistry.reader"
#   member     = "serviceAccount:${google_service_account.default["image-builder"].email}"
# }
