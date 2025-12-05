package devhelpermcp

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

const (
	devHelperName = "dev_helper"
)

func New(ctx context.Context, cfg Config) (tool.Toolset, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &dh{Config: cfg, name: devHelperName}, nil
}

type dh struct {
	Config

	name string
}

func (d *dh) Name() string {
	return d.name
}

func (d *dh) Tools(ctx agent.ReadonlyContext) ([]tool.Tool, error) {
	createCloudBuild, err := functiontool.New(
		functiontool.Config{
			Name:        "create_cloud_build",
			Description: "Starts a Google Cloud Build build. The source will be automatically cloned based on the tool parameters; there is no need to add a Cloud Build step to clone the source.",
		},
		d.createCloudBuildTool,
	)
	if err != nil {
		return nil, fmt.Errorf("createCloduBuild tool: %w", err)
	}

	getCloudBuild, err := functiontool.New(
		functiontool.Config{
			Name:        "get_cloud_build",
			Description: "Gets the details and status for a Google Cloud Build build.",
		},
		d.getCloudBuildTool,
	)
	if err != nil {
		return nil, fmt.Errorf("getCloudBuild tool: %w", err)
	}

	getCloudBuildLogs, err := functiontool.New(
		functiontool.Config{
			Name:        "get_cloud_build_logs",
			Description: "Gets the logs for a Google Cloud Build build.",
		},
		d.getCloudBuildLogsTool,
	)
	if err != nil {
		return nil, fmt.Errorf("getCloudBuildLogs tool: %w", err)
	}

	fetchTestOutput, err := functiontool.New(
		functiontool.Config{
			Name:        "fetch_test_output",
			Description: "Gets the test output when the test output logs have been uploaded to the test log repository.",
		},
		d.fetchTestOutputTool,
	)
	if err != nil {
		return nil, fmt.Errorf("fetchTestOutput tool: %w", err)
	}

	fetchProvenance, err := functiontool.New(
		functiontool.Config{
			Name:        "fetch_provenance",
			Description: "Gets the provenance for artifacts uploaded during a build.",
		},
		d.fetchProvenanceTool,
	)
	if err != nil {
		return nil, fmt.Errorf("fetchProvenance tool: %w", err)
	}

	createGoRepository, err := functiontool.New(
		functiontool.Config{
			Name:        "create_go_repository",
			Description: "Creates an ephemeral Go repository to hold go module artifacts. A random repository ID is generated and returned.",
		},
		d.createGoRepositoryTool,
	)

	return []tool.Tool{
		createCloudBuild,
		getCloudBuild,
		getCloudBuildLogs,
		fetchTestOutput,
		fetchProvenance,
		createGoRepository,
	}, nil
}

var _ tool.Toolset = (*dh)(nil)
