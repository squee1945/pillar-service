package service

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/google/go-github/v75/github"
)

//go:embed prompts/*.tmpl
var promptsFS embed.FS

func parsePromptTemplates(_ context.Context) (*template.Template, error) {
	return template.ParseFS(promptsFS, "prompts/*.tmpl")
}

func (s *Service) renderPrompt(ctx context.Context, pt promptTemplate) (string, error) {
	tname := pt.Name(ctx) + ".tmpl"

	data, err := pt.Data(ctx)
	if err != nil {
		return "", fmt.Errorf("getting prompt data: %w", err)
	}

	var buf bytes.Buffer
	if err := s.prompts.ExecuteTemplate(&buf, tname, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}
	return buf.String(), nil
}

type promptTemplate interface {
	Name(context.Context) string
	Data(context.Context) (any, error)
}

type promptReleasePublished struct {
	dependent string
	branch    string
	event     *github.ReleaseEvent
}

func (p *promptReleasePublished) Name(context.Context) string {
	return "release_published"
}

func (p *promptReleasePublished) Data(context.Context) (any, error) {
	js, err := eventJSON(p.event)
	if err != nil {
		return nil, err
	}
	return struct {
		Dependent string
		Branch    string
		EventJSON string
		Event     any
	}{
		Dependent: p.dependent,
		Branch:    p.branch,
		EventJSON: js,
		Event:     p.event,
	}, nil
}

func eventJSON(event any) (string, error) {
	js, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshalling JSON: %w", err)
	}
	return string(js), nil
}
