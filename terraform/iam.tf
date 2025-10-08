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
