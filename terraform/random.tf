resource "random_string" "suffix" {
  length  = 6     # Desired length for the short string
  special = false # Exclude special characters
  upper   = false
  lower   = true
  numeric = true
}
