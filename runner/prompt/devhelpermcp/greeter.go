package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type greeterInput struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type greeterOutput struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func greet(githubToken string) mcp.ToolHandlerFor[greeterInput, greeterOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input greeterInput) (*mcp.CallToolResult, greeterOutput, error) {
		return nil, greeterOutput{Greeting: "Hi " + input.Name}, nil
	}
}
