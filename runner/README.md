# Runner images

The images in the `runner` folder are used in Cloud Build pipelines.

## Building the prep image

$ gcloud builds submit \
  --project ${PROJECT_ID:?} \
  --region ${REGION:?} \
  runner/prep \
  --tag ${REGION:?}-docker.pkg.dev/${PROJECT_ID:?}/runner-images/prep:latest

## Building the prompt image

$ gcloud builds submit \
  --project ${PROJECT_ID:?} \
  --region ${REGION:?} \
  runner/prompt \
  --tag ${REGION:?}-docker.pkg.dev/${PROJECT_ID:?}/runner-images/prompt:latest
