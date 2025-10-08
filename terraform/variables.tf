variable "project_id" {
  description = "The project ID to host the resources in."
  type        = string
}

variable "region" {
  description = "The region to host the resources in."
  type        = string
  default     = "us-west1"
}

variable "github_app_id" {
  description = "The ID of the GitHub App."
  type        = number
}
