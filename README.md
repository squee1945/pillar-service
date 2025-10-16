# Pillar prototype

A Cloud Run service that listens to GitHub webhook events.
This service can fork a GitHub repo, then launch a "runner" (on Cloud Build)
that clones the forked repo and creates a development branch.

Next, it runs the Gemini CLI (in --yolo mode) with an arbitrary prompt.

This forms a general purpose event handling system.

## Create a GitHub app

TODO: Add instructions, but for now, you need to create your own and get the
appID for the next step.

## Installing

Most everything can be installed with `terraform`:

```
$ cd terraform
$ terraform apply -var 'project_id=<project>' -var 'github_app_id=<app_id>'
```

## Adding secrets

The terraform will create three secrets, but you need to manually add a
secret version with the actual secret.

  - `gemini-api-key`
  - `github-webhook-secret`
  - `github-private-key`

TODO: Improve these instructions.

## Update GitHub app

You must update the GitHub app to point the webhook handler to your
deployed Cloud Run service.

```
https://<cloud_run_service_name>.us-west1.run.app/webhook
```

TODO: Describe which webhook events to subscribe to. Importantly, you must
subscribe to release events for the current example.

## Install the GitHub app

Install the GitHub app on a repository.

TODO: Improve these instructions.

## Generate a webhook event

Currently, only the creation of a release triggers handling, so on the
repository where you installed the GitHub app, create a new release.

You should see the webhook event hit your Cloud Run logs, then see a runner
started in Cloud Build.
