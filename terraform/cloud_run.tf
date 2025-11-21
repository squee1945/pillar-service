resource "google_cloud_run_v2_service" "default" {
  project  = var.project_id
  name     = "pillar-service"
  location = var.region

  template {
    service_account = google_service_account.default["pillar-service"].email
    scaling {
      max_instance_count = 10
    }
    containers {
      image = ko_build.pillar_service.image_ref
      ports {
        container_port = 8080
      }
      env {
        name  = "PROJECT_ID"
        value = var.project_id
      }
      env {
        name  = "REGION"
        value = var.region
      }
      env {
        name  = "KMS_KEY_NAME"
        value = google_kms_crypto_key.default.id
      }
      env {
        name  = "RUNNER_SERVICE_ACCOUNT"
        value = google_service_account.default["runner"].id
      }
      env {
        name  = "SUB_BUILD_SERVICE_ACCOUNT"
        value = google_service_account.default["sub-build"].id
      }
      env {
        name  = "SUB_BUILD_LOGS_BUCKET"
        value = google_storage_bucket.sub_build_logs.id
      }
      env {
        name  = "PREP_IMAGE"
        value = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.runner_images.repository_id}/prep:latest"
      }
      env {
        name  = "PROMPT_IMAGE"
        value = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.runner_images.repository_id}/prompt:latest"
      }
      env {
        name  = "PROMPT_BUCKET"
        value = google_storage_bucket.prompt_bucket.name
      }
      env {
        name  = "GITHUB_APP_ID"
        value = var.github_app_id
      }
      env {
        name  = "GITHUB_WEBHOOK_SECRET_NAME"
        value = "${google_secret_manager_secret.default["github-webhook-secret"].name}/versions/latest"
      }
      env {
        name  = "GITHUB_PRIVATE_KEY_SECRET_NAME"
        value = "${google_secret_manager_secret.default["github-private-key"].name}/versions/latest"
      }
      env {
        name  = "GEMINI_API_KEY_SECRET_NAME"
        value = "${google_secret_manager_secret.default["gemini-api-key"].name}/versions/latest"
      }
    }
  }
}

resource "google_cloud_run_v2_service_iam_member" "noauth" {
  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.default.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
