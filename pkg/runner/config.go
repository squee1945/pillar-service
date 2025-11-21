package runner

import (
	"fmt"
	"time"

	"github.com/squee1945/pillar-service/pkg/logger"
)

type Config struct {
	Log            logger.L
	ProjectID      string
	Region         string
	PromptBucket   string
	KMSKeyName     string
	ServiceAccount string

	PrepImage     string
	GitHubToken   string
	Owner         string
	Repo          string
	DefaultBranch string
	DevBranch     string

	PromptImage  string
	GeminiAPIKey string
	Prompt       string // If empty, Gemini CLI not invoked.

	SubBuildServiceAccount string
	SubBuildLogsBucket     string

	// Optional config
	RunnerTimeout         time.Duration
	GeminiMaxSessionTurns int
	DevHelperIncludeTools []string
	DevHelperExcludeTools []string
	DevHelperMCPTimeout   time.Duration
	GithubIncludeTools    []string
	GithubExcludeTools    []string
	GithubMCPTimeout      time.Duration
}

func (c Config) validate() error {
	if c.ProjectID == "" {
		return fmt.Errorf("ProjectID must be set")
	}
	if c.Region == "" {
		return fmt.Errorf("Region must be set")
	}
	if c.PromptBucket == "" {
		return fmt.Errorf("PromptBucket must be set")
	}
	if c.KMSKeyName == "" {
		return fmt.Errorf("KMSKeyName must be set")
	}
	if c.ServiceAccount == "" {
		return fmt.Errorf("ServiceAccount must be set")
	}
	if c.PrepImage == "" {
		return fmt.Errorf("PrepImage must be set")
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
	if c.DevBranch == "" {
		return fmt.Errorf("DevBranch must be set")
	}
	if c.PromptImage == "" {
		return fmt.Errorf("PromptImage must be set")
	}
	if c.SubBuildServiceAccount == "" {
		return fmt.Errorf("SubBuildServiceAccount must be set")
	}
	if c.SubBuildLogsBucket == "" {
		return fmt.Errorf("SubBuildLogsBucket must be set")
	}
	return nil
}
