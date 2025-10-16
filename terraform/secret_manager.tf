resource "google_secret_manager_secret" "default" {
  project = var.project_id

  for_each = toset([
    "gemini-api-key",
    "github-webhook-secret",
    "github-private-key"
  ])

  secret_id = each.key

  replication {
    user_managed {
      replicas {
        location = var.region
      }
    }
  }

  depends_on = [
    google_project_service.default
  ]
}
