package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/go-github/v75/github"
	"github.com/squee1945/pillar-service/pkg/runner"
)

func (s *Service) run(ctx context.Context, installationID int64, repo *github.Repository, prompt string) error {
	cfg, err := s.runnerConfig(ctx, installationID, repo)
	if err != nil {
		return fmt.Errorf("generating runner config: %v", err)
	}

	cfg.Prompt = prompt

	r, err := runner.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("creating runner: %v", err)
	}
	return r.Run(ctx)
}

func (s *Service) runnerConfigBase(ctx context.Context) (runner.Config, error) {
	geminiAPIKey, err := s.Secrets.Read(ctx, s.GeminiAPIKeySecretName)
	if err != nil {
		return runner.Config{}, fmt.Errorf("getting Gemini API key: %v", err)
	}

	return runner.Config{
		Log:                   s.Log,
		ProjectID:             s.ProjectID,
		Region:                s.Region,
		PromptBucket:          s.PromptBucket,
		KMSKeyName:            s.KMSKeyName,
		ServiceAccount:        s.RunnerServiceAccount,
		PrepImage:             s.PrepImage,
		DevBranch:             fmt.Sprintf("%s-%d-%s", s.ServiceName, time.Now().Unix(), randomString(4)),
		PromptImage:           s.PromptImage,
		RunnerTimeout:         s.RunnerTimeout,
		MCPToolTimeout:        s.MCPToolTimeout,
		GeminiMaxSessionTurns: s.GeminiMaxSessionTurns,
		GeminiAPIKey:          string(geminiAPIKey),
	}, nil
}

func (s *Service) runnerConfig(ctx context.Context, installationID int64, repo *github.Repository) (runner.Config, error) {
	githubToken, err := s.installationToken(ctx, installationID, withRepoIDs(repo.GetID()))
	if err != nil {
		return runner.Config{}, fmt.Errorf("generating installation token: %v", err)
	}

	cfg, err := s.runnerConfigBase(ctx)
	if err != nil {
		return runner.Config{}, err
	}
	cfg.GitHubToken = githubToken
	cfg.Owner = repo.GetOwner().GetLogin()
	cfg.Repo = repo.GetName()
	cfg.DefaultBranch = repo.GetDefaultBranch()
	return cfg, nil

}

const consonants = "bcdfghjklmnpqrstvwxyz"

func randomString(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	charsetLength := len(consonants)
	for i := range b {
		b[i] = consonants[rand.Intn(charsetLength)]
	}
	return string(b)
}
