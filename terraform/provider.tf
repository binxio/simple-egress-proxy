provider "google" {
  region = "europe-west4"
}

data "google_client_config" "current" {
}
