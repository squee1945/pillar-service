resource "google_storage_bucket" "prompt_bucket" {
  project                     = var.project_id
  name                        = "prompt-bucket-${var.project_id}-${random_string.suffix.result}"
  location                    = var.region
  force_destroy               = true
  uniform_bucket_level_access = true

  lifecycle_rule {
    condition {
      age = 1
    }
    action {
      type = "Delete"
    }
  }
}
