package adkrunner

import (
	"fmt"
	"time"

	"github.com/squee1945/pillar-service/pkg/logger"
)

type Config struct {
	Log          logger.L
	ServiceName  string
	ProjectID    string
	Region       string
	GeminiAPIKey string
	Prompt       string

	GitHubToken   string
	Owner         string
	Repo          string
	DefaultBranch string

	SubBuildServiceAccount   string
	SubBuildLogsBucket       string
	SubBuildTestOutputBucket string

	// Optional config
	RunnerTimeout         time.Duration
	GeminiMaxSessionTurns int
	DevHelperIncludeTools []string
	DevHelperExcludeTools []string
	GithubIncludeTools    []string
	GithubExcludeTools    []string
}

func (c Config) validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("ServiceName must be set")
	}
	if c.ProjectID == "" {
		return fmt.Errorf("ProjectID must be set")
	}
	if c.Region == "" {
		return fmt.Errorf("Region must be set")
	}
	if c.GeminiAPIKey == "" {
		return fmt.Errorf("GeminiAPIKey must be set")
	}
	if c.Prompt == "" {
		return fmt.Errorf("Prompt must be set")
	}
	if c.GitHubToken == "" {
		return fmt.Errorf("GitHubToken must be set")
	}
	if c.Owner == "" {
		return fmt.Errorf("Owner must be set")
	}
	if c.Repo == "" {
		return fmt.Errorf("Repo must be set")
	}
	if c.DefaultBranch == "" {
		return fmt.Errorf("DefaultBranch must be set")
	}
	if c.SubBuildServiceAccount == "" {
		return fmt.Errorf("SubBuildServiceAccount must be set")
	}
	if c.SubBuildLogsBucket == "" {
		return fmt.Errorf("SubBuildLogsBucket must be set")
	}
	if c.SubBuildTestOutputBucket == "" {
		return fmt.Errorf("SubBuildTestOutputBucket must be set")
	}
	return nil
}
