resource "google_project_service" "default" {
  project = var.project_id

  for_each = toset([
    "cloudbuild.googleapis.com",
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudkms.googleapis.com",
    "secretmanager.googleapis.com",
  ])

  service = each.key
}
