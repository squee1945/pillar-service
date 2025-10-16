resource "null_resource" "build_runner_images" {
  for_each = toset([
    "prep",
    "prompt",
  ])

  triggers = {
    source_code_hash = jsonencode([
      for f in fileset("../runner/${each.key}", "**") : filebase64sha256("../runner/${each.key}/${f}")
    ])
    daily_rebuild = time_rotating.daily.rfc3339
  }

  provisioner "local-exec" {
    command = join(" ", [
      "gcloud",
      "builds",
      "submit",
      "--project=${var.project_id}",
      "--region=${var.region}",
      "--tag=${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.runner_images.repository_id}/${each.key}:latest",
      "../runner/${each.key}",
    ])
  }
}

resource "ko_build" "pillar_service" {
  importpath = "github.com/squee1945/pillar-service/cmd/web"
  repo       = "${google_artifact_registry_repository.pillar_service.location}-docker.pkg.dev/${google_artifact_registry_repository.pillar_service.project}/${google_artifact_registry_repository.pillar_service.repository_id}/web"
  env        = ["CGO_ENABLED=0"]
  base_image = "gcr.io/distroless/static-debian12"
}
