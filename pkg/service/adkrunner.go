package service

import (
	"context"
	"fmt"

	"github.com/google/go-github/v75/github"
	"github.com/squee1945/pillar-service/pkg/adkrunner"
)

func (s *Service) runADK(ctx context.Context, installationID int64, repo *github.Repository, prompt string, configOpts ...adkConfigOption) error {
	cfg, err := s.runnerADKConfig(ctx, installationID, repo)
	if err != nil {
		return fmt.Errorf("generating runner config: %v", err)
	}

	for _, o := range configOpts {
		o(&cfg)
	}

	cfg.Prompt = prompt

	r, err := adkrunner.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("creating runner: %v", err)
	}
	return r.Run(ctx)
}

func (s *Service) runnerADKConfigBase(ctx context.Context) (adkrunner.Config, error) {
	geminiAPIKey, err := s.Secrets.Read(ctx, s.GeminiAPIKeySecretName)
	if err != nil {
		return adkrunner.Config{}, fmt.Errorf("getting Gemini API key: %v", err)
	}

	return adkrunner.Config{
		Log:                      s.Log,
		ServiceName:              s.ServiceName,
		ProjectID:                s.ProjectID,
		Region:                   s.Region,
		GeminiAPIKey:             string(geminiAPIKey),
		SubBuildServiceAccount:   s.SubBuildServiceAccount,
		SubBuildLogsBucket:       s.SubBuildLogsBucket,
		SubBuildTestOutputBucket: s.SubBuildTestOutputBucket,
	}, nil
}

func (s *Service) runnerADKConfig(ctx context.Context, installationID int64, repo *github.Repository) (adkrunner.Config, error) {
	githubToken, err := s.installationToken(ctx, installationID, withRepoIDs(repo.GetID()))
	if err != nil {
		return adkrunner.Config{}, fmt.Errorf("generating installation token: %v", err)
	}

	cfg, err := s.runnerADKConfigBase(ctx)
	if err != nil {
		return adkrunner.Config{}, err
	}
	cfg.GitHubToken = githubToken
	cfg.Owner = repo.GetOwner().GetLogin()
	cfg.Repo = repo.GetName()
	cfg.DefaultBranch = repo.GetDefaultBranch()
	return cfg, nil
}

type adkConfigOption func(*adkrunner.Config)

func withDevHelperIncludeToolsADK(tools ...string) adkConfigOption {
	return func(cfg *adkrunner.Config) {
		cfg.DevHelperIncludeTools = tools
	}
}

func withDevHelperExcludeToolsADK(tools ...string) adkConfigOption {
	return func(cfg *adkrunner.Config) {
		cfg.DevHelperExcludeTools = tools
	}
}

func withGithubIncludeToolsADK(tools ...string) adkConfigOption {
	return func(cfg *adkrunner.Config) {
		cfg.GithubIncludeTools = tools
	}
}

func withGithubExcludeToolsADK(tools ...string) adkConfigOption {
	return func(cfg *adkrunner.Config) {
		cfg.GithubExcludeTools = tools
	}
}
