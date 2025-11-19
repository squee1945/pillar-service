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
	githubToken = flag.String("github_token", "", "GitHub access token")
)

func main() {
	flag.Parse()

	if *githubToken == "" {
		log.Fatal("github_token is required")
	}

	i := &mcp.Implementation{
		Name:    "dev_helper",
		Title:   "Developer Helper - High level tools to assist in creating contributions to GitHub repositories.",
		Version: "v1.0.0",
	}
	opts := &mcp.ServerOptions{HasTools: true}
	server := mcp.NewServer(i, opts)

	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, greeterTool(*githubToken))
	mcp.AddTool(server, &mcp.Tool{Name: "prep_dev_env", Description: "Prepares a dev environment to facilitate a contribution against an upstream repository."}, prepDevEnvTool(*githubToken))

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
