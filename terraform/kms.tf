resource "google_kms_key_ring" "default" {
  project  = var.project_id
  name     = "default"
  location = var.region

  depends_on = [
    google_project_service.default
  ]
}

resource "google_kms_crypto_key" "default" {
  name     = "default"
  key_ring = google_kms_key_ring.default.id
  purpose  = "ENCRYPT_DECRYPT"

  version_template {
    algorithm = "GOOGLE_SYMMETRIC_ENCRYPTION"
  }
}
