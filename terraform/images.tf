# resource "null_resource" "build_pillar_service_image" {
#   triggers = {
#     always_run = timestamp()
#   }

#   provisioner "local-exec" {
#     command = "gcloud --project=${var.project_id} builds submit --config ../cloudbuild.yaml --substitutions=_SERVICE_IMAGE=${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.pillar_service.repository_id}/pillar-service .."
#   }
# }

resource "ko_build" "pillar_service" {
  importpath = "github.com/squee1945/pillar-service/cmd/web"
  repo       = "${google_artifact_registry_repository.pillar_service.location}-docker.pkg.dev/${google_artifact_registry_repository.pillar_service.project}/${google_artifact_registry_repository.pillar_service.repository_id}/web"
  env        = ["CGO_ENABLED=0"]
  base_image = "gcr.io/distroless/static-debian12"
}
