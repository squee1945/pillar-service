data "google_project" "project" {}

resource "google_service_account" "default" {
  project = var.project_id

  for_each = {
    "pillar-service" = "Service Account for the Pillar Service"
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

resource "google_project_iam_member" "pillar_service_cloudbuild_editor" {
  project = var.project_id
  role    = "roles/cloudbuild.builds.editor"
  member  = "serviceAccount:${google_service_account.default["pillar-service"].email}"
}

resource "google_service_account_iam_member" "pillar_service_can_impersonate_sub_build_sa" {
  service_account_id = google_service_account.default["sub-build"].name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.default["pillar-service"].email}"
}

resource "google_storage_bucket_iam_member" "pillar_service_can_read_build_logs" {
  bucket = google_storage_bucket.sub_build_logs.name
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${google_service_account.default["pillar-service"].email}"
}

resource "google_storage_bucket_iam_member" "pillar_service_can_read_test_output" {
  bucket = google_storage_bucket.sub_build_test_output.name
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${google_service_account.default["pillar-service"].email}"
}

resource "google_project_iam_member" "pillar_service_can_create_repositories" {
  project = var.project_id
  role    = "roles/artifactregistry.admin"
  member  = "serviceAccount:${google_service_account.default["pillar-service"].email}"
}

resource "google_project_iam_member" "pillar_service_can_fetch_provenance" {
  project = var.project_id
  role    = "roles/containeranalysis.occurrences.viewer"
  member  = "serviceAccount:${google_service_account.default["pillar-service"].email}"
}

resource "google_storage_bucket_iam_member" "sub_build_log_bucket_viewer" {
  bucket = google_storage_bucket.sub_build_logs.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.default["sub-build"].email}"
}

resource "google_project_iam_member" "sub_build_can_write_artifacts" {
  project = var.project_id
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:${google_service_account.default["sub-build"].email}"
}

resource "google_storage_bucket_iam_member" "sub_build_test_output_creator" {
  bucket = google_storage_bucket.sub_build_test_output.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.default["sub-build"].email}"
}
