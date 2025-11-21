package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	cloudbuildpb "cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
)

const (
	maxLogs = 10 * 1024 * 1024
)

type createCloudBuildInput struct {
	CloudBuildYAML string `json:"cloud_build_yaml" jsonschema:"Serialized YAML for the cloudbuild.yaml.`
	Owner          string `json:"source_owner" jsonschema:"The owner of the source repo."`
	Repo           string `json:"source_repo" jsonschema:"The name of the source repo.`
	Commit         string `json:"source_commit" jsonschema:"The commit sha to clone the repo at."`
}

func (i createCloudBuildInput) validate() error {
	var errs []error
	if i.CloudBuildYAML == "" {
		errs = append(errs, errors.New("cloud_build_yaml is required."))
	}
	if i.Owner == "" {
		errs = append(errs, errors.New("source_owner is required."))
	}
	if i.Repo == "" {
		errs = append(errs, errors.New("source_repo is required."))
	}
	if i.Commit == "" {
		errs = append(errs, errors.New("source_commit is required."))
	}
	return errors.Join(errs...)
}

type createCloudBuildOutput struct {
	BuildID string `json:"build_id" jsonschema:"The Build ID of the created build.`
}

func createCloudBuildTool(githubToken, projectID, region, subBuildServiceAccount, subBuildLogsBucket string) mcp.ToolHandlerFor[createCloudBuildInput, createCloudBuildOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input createCloudBuildInput) (*mcp.CallToolResult, createCloudBuildOutput, error) {
		if err := input.validate(); err != nil {
			return nil, createCloudBuildOutput{}, err
		}

		var build cloudbuildpb.Build
		if err := yaml.Unmarshal([]byte(input.CloudBuildYAML), &build); err != nil {
			return nil, createCloudBuildOutput{}, fmt.Errorf("yaml.Unmarshal: failed to parse YAML config: %w", err)
		}

		if build.Options == nil {
			build.Options = &cloudbuildpb.BuildOptions{}
		}
		build.Options.LogStreamingOption = cloudbuildpb.BuildOptions_STREAM_ON
		build.Options.Logging = cloudbuildpb.BuildOptions_GCS_ONLY
		build.LogsBucket = "gs://" + subBuildLogsBucket
		build.ServiceAccount = subBuildServiceAccount
		build.Source = &cloudbuildpb.Source{
			Source: &cloudbuildpb.Source_GitSource{
				GitSource: &cloudbuildpb.GitSource{
					Url:      fmt.Sprintf("https://github.com/%s/%s", input.Owner, input.Repo),
					Revision: input.Commit,
				},
			},
		}

		client, err := cloudBuildClient(ctx, region)
		if err != nil {
			return nil, createCloudBuildOutput{}, err
		}
		defer client.Close()

		createReq := &cloudbuildpb.CreateBuildRequest{
			Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, region),
			Build:  &build,
		}

		op, err := client.CreateBuild(ctx, createReq)
		if err != nil {
			return nil, createCloudBuildOutput{}, fmt.Errorf("client.CreateBuild failed: %w", err)
		}

		metadata, err := op.Metadata()
		if err != nil {
			return nil, createCloudBuildOutput{}, fmt.Errorf("failed to get operation metadata: %w", err)
		}

		buildID := metadata.GetBuild().GetId()

		return nil, createCloudBuildOutput{BuildID: buildID}, nil
	}
}

type getCloudBuildInput struct {
	BuildID string `json:"build_id" jsonschema:"The Build ID of the build.`
}

func (i getCloudBuildInput) validate() error {
	var errs []error
	if i.BuildID == "" {
		errs = append(errs, errors.New("build_id is required."))
	}
	return errors.Join(errs...)
}

type getCloudBuildOutput struct {
	BuildID   string `json:"build_id" jsonschema:"The Build ID of the build.`
	Status    string `json:"status" jsonschema:"The build status.`
	BuildJSON string `json:"build_json" jsonschema:"The build details, as serialized JSON`
}

func getCloudBuildTool(projectID, region string) mcp.ToolHandlerFor[getCloudBuildInput, getCloudBuildOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input getCloudBuildInput) (*mcp.CallToolResult, getCloudBuildOutput, error) {
		if err := input.validate(); err != nil {
			return nil, getCloudBuildOutput{}, err
		}

		build, err := getCloudBuild(ctx, projectID, region, input.BuildID)
		if err != nil {
			return nil, getCloudBuildOutput{}, err
		}

		buildJSON, err := json.MarshalIndent(build, "", "  ")
		if err != nil {
			return nil, getCloudBuildOutput{}, fmt.Errorf("json.MarshalIndent: %w", err)
		}

		output := getCloudBuildOutput{
			BuildID:   input.BuildID,
			Status:    build.GetStatus().String(),
			BuildJSON: string(buildJSON),
		}

		return nil, output, nil
	}
}

type getCloudBuildLogsInput struct {
	BuildID string `json:"build_id" jsonschema:"The Build ID of the build.`
}

func (i getCloudBuildLogsInput) validate() error {
	var errs []error
	if i.BuildID == "" {
		errs = append(errs, errors.New("build_id is required."))
	}
	return errors.Join(errs...)
}

type getCloudBuildLogsOutput struct {
	Logs []string `json:"logs" jsonschema:"The log lines."`
}

func getCloudBuildLogsTool(projectID, region string) mcp.ToolHandlerFor[getCloudBuildLogsInput, getCloudBuildLogsOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input getCloudBuildLogsInput) (*mcp.CallToolResult, getCloudBuildLogsOutput, error) {
		if err := input.validate(); err != nil {
			return nil, getCloudBuildLogsOutput{}, err
		}

		build, err := getCloudBuild(ctx, projectID, region, input.BuildID)
		if err != nil {
			return nil, getCloudBuildLogsOutput{}, err
		}

		storageClient, err := storage.NewClient(ctx)
		if err != nil {
			return nil, getCloudBuildLogsOutput{}, fmt.Errorf("storage.NewClient: %w", err)
		}
		defer storageClient.Close()

		bucket := strings.TrimPrefix(build.LogsBucket, "gs://")
		object := fmt.Sprintf("log-%s.txt", input.BuildID)

		rc, err := storageClient.Bucket(bucket).Object(object).NewReader(ctx)
		if err != nil {
			return nil, getCloudBuildLogsOutput{}, fmt.Errorf("getting logs object (bucket %q object %q): %w", bucket, object, err)
		}
		defer rc.Close()

		logBlob, err := io.ReadAll(io.LimitReader(rc, maxLogs))
		lines := strings.Split(string(logBlob), "\n")

		output := getCloudBuildLogsOutput{
			Logs: lines,
		}

		return nil, output, nil
	}
}

func cloudBuildClient(ctx context.Context, region string) (*cloudbuild.Client, error) {
	endpoint := fmt.Sprintf("%s-cloudbuild.googleapis.com:443", region)
	return cloudbuild.NewClient(ctx, option.WithEndpoint(endpoint))
}

func getCloudBuild(ctx context.Context, projectID, region, buildID string) (*cloudbuildpb.Build, error) {
	client, err := cloudBuildClient(ctx, region)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	gbReq := &cloudbuildpb.GetBuildRequest{
		Name:      fmt.Sprintf("projects/%s/locations/%s/builds/%s", projectID, region, buildID),
		ProjectId: projectID,
		Id:        buildID,
	}

	return client.GetBuild(ctx, gbReq)
}
