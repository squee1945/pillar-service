terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.6.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "5.15.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.5.1"
    }
    ko = {
      source  = "ko-build/ko"
      version = "0.0.17"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.11.1"
    }
  }
}

provider "google" {
  project = var.project_id
}

provider "google-beta" {
  project = var.project_id
}
