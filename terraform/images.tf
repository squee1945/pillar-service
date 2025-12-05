resource "ko_build" "pillar_service" {
  importpath = "github.com/squee1945/pillar-service/cmd/web"
  repo       = "${google_artifact_registry_repository.pillar_service.location}-docker.pkg.dev/${google_artifact_registry_repository.pillar_service.project}/${google_artifact_registry_repository.pillar_service.repository_id}/web"
  env        = ["CGO_ENABLED=0"]
  base_image = "gcr.io/distroless/static-debian12"
}
