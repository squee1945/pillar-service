// This package is a local MCP server with a set of tools purpose-built
// to support the Pillar Gemini agent.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	githubToken            = flag.String("github_token", "", "GitHub access token")
	subBuildServiceAccount = flag.String("sub_build_service_account", "", "Service account for sub-build")
	subBuildLogsBucket     = flag.String("sub_build_logs_bucket", "", "Log bucker for sub-build")
	projectID              = flag.String("project_id", "", "The project ID")
	region                 = flag.String("region", "", "The region")
)

func main() {
	flag.Parse()

	if *githubToken == "" {
		log.Fatal("--github_token is required")
	}

	if *subBuildServiceAccount == "" {
		log.Fatal("--sub_build_service_account is required")
	}

	if *subBuildLogsBucket == "" {
		log.Fatal("--sub_build_logs_bucket is required")
	}

	if *projectID == "" {
		log.Fatal("--project_id is required")
	}

	if *region == "" {
		log.Fatal("--region is required")
	}

	i := &mcp.Implementation{
		Name:    "dev_helper",
		Title:   "Developer Helper - High level tools to assist in creating contributions to GitHub repositories.",
		Version: "v1.0.0",
	}
	opts := &mcp.ServerOptions{HasTools: true}
	server := mcp.NewServer(i, opts)

	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, greeterTool(*githubToken))

	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "prep_dev_env",
			Description: "Prepares a dev environment to facilitate a contribution against an upstream repository.",
		},
		prepDevEnvTool(*githubToken),
	)

	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "create_cloud_build",
			Description: "Starts a Google Cloud Build build. The source will be automatically cloned based on the tool parameters; there is no need to add a Cloud Build step to clone the source.",
		},
		createCloudBuildTool(*githubToken, *projectID, *region, *subBuildServiceAccount, *subBuildLogsBucket),
	)

	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "get_cloud_build",
			Description: "Gets the details and status for a Google Cloud Build build.",
		},
		getCloudBuildTool(*projectID, *region),
	)

	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "get_cloud_build_logs",
			Description: "Gets the logs for a Google Cloud Build build.",
		},
		getCloudBuildLogsTool(*projectID, *region),
	)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
