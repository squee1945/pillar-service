data "google_project" "project" {}

resource "google_service_account" "default" {
  project = var.project_id

  for_each = {
    "pillar-service" = "Service Account for the Pillar Service"
    "runner"         = "Service Account for the Cloud Build runner"
    "sub-build"      = "Service Account for builds created by the agent."
  }

  account_id   = each.key
  display_name = each.value
}

resource "google_project_iam_member" "pillar_service_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.default["pillar-service"].email}"
}

resource "google_project_iam_member" "pillar_service_kms_encryptor" {
  project = var.project_id
  role    = "roles/cloudkms.cryptoKeyEncrypter"
  member  = "serviceAccount:${google_service_account.default["pillar-service"].email}"
}

resource "google_storage_bucket_iam_member" "pillar_service_gcs_writer" {
  bucket = google_storage_bucket.prompt_bucket.name
  role   = "roles/storage.objectCreator"
  member = "serviceAccount:${google_service_account.default["pillar-service"].email}"
}

resource "google_project_iam_member" "pillar_service_cloudbuild_editor" {
  project = var.project_id
  role    = "roles/cloudbuild.builds.editor"
  member  = "serviceAccount:${google_service_account.default["pillar-service"].email}"
}

resource "google_service_account_iam_member" "pillar_service_can_impersonate_cli_runner" {
  service_account_id = google_service_account.default["runner"].name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.default["pillar-service"].email}"
}

resource "google_project_iam_member" "runner_kms_decryptor" {
  project = var.project_id
  role    = "roles/cloudkms.cryptoKeyDecrypter"
  member  = "serviceAccount:${google_service_account.default["runner"].email}"
}

resource "google_artifact_registry_repository_iam_member" "runner_image_reader" {
  project  = var.project_id
  provider = google-beta

  repository = google_artifact_registry_repository.runner_images.id
  role       = "roles/artifactregistry.reader"
  member     = "serviceAccount:${google_service_account.default["runner"].email}"
}

resource "google_project_iam_member" "runner_logs_writer" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.default["runner"].email}"
}

resource "google_project_iam_member" "runner_storage_viewer" {
  project = var.project_id
  role    = "roles/storage.objectViewer"
  member  = "serviceAccount:${google_service_account.default["runner"].email}"
}

resource "google_project_iam_member" "runner_cloud_build_creator" {
  project = var.project_id
  role    = "roles/cloudbuild.builds.editor"
  member  = "serviceAccount:${google_service_account.default["runner"].email}"
}

resource "google_storage_bucket_iam_member" "runner_logs_reader" {
  bucket = google_storage_bucket.sub_build_logs.name
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${google_service_account.default["runner"].email}"
}

resource "google_project_iam_member" "runner_container_analysis_viewer" {
  project = var.project_id
  role    = "roles/containeranalysis.occurrences.viewer"
  member  = "serviceAccount:${google_service_account.default["runner"].email}"
}

resource "google_service_account_iam_member" "cli_runner_can_impersonate_sub_build" {
  service_account_id = google_service_account.default["sub-build"].name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.default["runner"].email}"
}

resource "google_storage_bucket_iam_member" "sub_build_log_bucket_viewer" {
  bucket = google_storage_bucket.sub_build_logs.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.default["sub-build"].email}"
}

resource "google_artifact_registry_repository_iam_member" "go_repo_writer" {
  repository = google_artifact_registry_repository.sub_build_go_repository.id
  role       = "roles/artifactregistry.writer"
  member     = "serviceAccount:${google_service_account.default["sub-build"].email}"
}

resource "google_storage_bucket_iam_member" "sub_build_test_output_creator" {
  bucket = google_storage_bucket.sub_build_test_output.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.default["sub-build"].email}"
}
