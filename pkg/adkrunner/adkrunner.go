// This package uses the ADK library to run the agentic loop.
package adkrunner

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/squee1945/pillar-service/pkg/adkrunner/devhelpermcp"
	"github.com/squee1945/pillar-service/pkg/adkrunner/githubmcp"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
	"google.golang.org/protobuf/proto"
)

const (
	cols = 160

	// Thinking budget is TODO-magic.
	// https://docs.cloud.google.com/vertex-ai/generative-ai/docs/thinking
	thinkingBudget = 4096

	// TopP is the top-p used for all models.
	// https://cloud.google.com/vertex-ai/generative-ai/docs/multimodal/content-generation-parameters#top-p
	topP = 0.5

	// Temperature is the temperature used for all models.
	// https://cloud.google.com/vertex-ai/generative-ai/docs/multimodal/content-generation-parameters#temperature
	temperature = 0.35

	geminiVersion = "gemini-3-pro-preview" // "gemini-2.5-pro"
)

type R struct {
	Config
}

func New(ctx context.Context, cfg Config) (*R, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &R{Config: cfg}, nil
}

func (r *R) Run(ctx context.Context) error {
	model, err := gemini.NewModel(ctx, geminiVersion, &genai.ClientConfig{APIKey: r.GeminiAPIKey})
	if err != nil {
		return fmt.Errorf("gemini.NewModel(): %w", err)
	}

	githubTools, err := githubmcp.New(ctx, r.GitHubToken, r.GithubIncludeTools, r.GithubExcludeTools)
	if err != nil {
		return fmt.Errorf("githubmcp.New(): %w", err)
	}

	cfg := devhelpermcp.Config{
		ServiceName:              r.ServiceName,
		GitHubToken:              r.GitHubToken,
		ProjectID:                r.ProjectID,
		Region:                   r.Region,
		SubBuildServiceAccount:   r.SubBuildServiceAccount,
		SubBuildLogsBucket:       r.SubBuildLogsBucket,
		SubBuildTestOutputBucket: r.SubBuildTestOutputBucket,
		IncludedTools:            r.DevHelperIncludeTools,
		ExcludedTools:            r.DevHelperExcludeTools,
	}

	devHelperTools, err := devhelpermcp.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("devhelpermcp.New(): %w", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "pillar_agent",
		Model:       model,
		Description: "Agent to assist in OSS maintenance tasks.",
		// InstructionProvider: func(ctx agent.ReadonlyContext) (string, error) { return r.Prompt, nil },
		Toolsets: []tool.Toolset{
			githubTools,
			devHelperTools,
		},
		GenerateContentConfig: &genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: proto.Int32(thinkingBudget),
			},
			Temperature: proto.Float32(temperature),
			TopP:        proto.Float32(topP),
		},
		BeforeToolCallbacks: []llmagent.BeforeToolCallback{
			func(ctx tool.Context, tool tool.Tool, args map[string]any) (map[string]any, error) {
				r.Log.Info(ctx, "TOOL CALL %s(%q)", tool.Name(), args)
				return nil, nil
			},
		},
		AfterToolCallbacks: []llmagent.AfterToolCallback{
			func(ctx tool.Context, tool tool.Tool, args, result map[string]any, err error) (map[string]any, error) {
				if err != nil {
					r.Log.Error(ctx, "TOOL ERR  %s(%q) = %q: %v", tool.Name(), args, result, err)
				} else {
					r.Log.Info(ctx, "TOOL OKAY %s(%q) = %q", tool.Name(), args, result)
				}
				return nil, nil
			},
		},
	})
	if err != nil {
		return err
	}

	const app = "TODO-app"
	const userID = "TODO-userID"

	service := session.InMemoryService()

	resp, err := service.Create(ctx, &session.CreateRequest{
		AppName: app,
		UserID:  userID,
		// State:   map[string]any{"BUILD_ID": "{BUILD_ID}"},
	})
	if err != nil {
		return fmt.Errorf("service.Create(): %v", err)
	}
	sessionID := resp.Session.ID()

	rn, err := runner.New(runner.Config{
		AppName:        app,
		Agent:          a,
		SessionService: service,
	})
	if err != nil {
		return fmt.Errorf("creating ADK runner: %v", err)
	}

	msg := genai.NewContentFromText(r.Prompt, genai.RoleUser)

	it := rn.Run(ctx, userID, sessionID, msg, agent.RunConfig{
		StreamingMode: agent.StreamingModeNone,
	})

	var citations []*genai.Citation
	for event, err := range it {
		if err != nil {
			return err
		}
		if event == nil {
			r.Log.Warn(ctx, "ADK response is empty: %v", event)
			continue
		}
		llmResp := event.LLMResponse
		if llmResp.ErrorCode != "" || llmResp.ErrorMessage != "" {
			return fmt.Errorf("ADK LLM response error: %v: %v", llmResp.ErrorCode, llmResp.ErrorMessage)
		}

		if llmResp.Content == nil {
			r.Log.Warn(ctx, "ADK response content is nil: %v", llmResp)
			continue
		}

		for _, p := range llmResp.Content.Parts {
			if s := strings.TrimSpace(p.Text); s != "" {
				for _, line := range wrap(s, cols, ">>> ") {
					r.Log.Debug(ctx, line)
				}
			}
		}
		if cm := llmResp.CitationMetadata; cm != nil {
			for _, c := range cm.Citations {
				if c.Title == "" || c.License == "" || c.PublicationDate == (civil.Date{}) || c.URI == "" {
					continue
				}
				citations = append(citations, c)
			}
		}
	}

	return nil
}
