data "google_project" "project" {}

resource "google_service_account" "default" {
  project = var.project_id

  for_each = {
    "pillar-service" = "Service Account for the Pillar Service"
    # "image-builder"       = "Service Account for building the image from runner-image/src and deploying to AR"
    # "cli-runner"          = "Service Account for running the Gemini CLI in Cloud Build"
  }

  account_id   = each.key
  display_name = each.value
}

resource "google_project_iam_member" "pillar_service_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.default["pillar-service"].email}"
}

