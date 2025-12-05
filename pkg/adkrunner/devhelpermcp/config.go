package devhelpermcp

import "errors"

type Config struct {
	ServiceName              string
	GitHubToken              string
	ProjectID                string
	Region                   string
	SubBuildServiceAccount   string
	SubBuildLogsBucket       string
	SubBuildTestOutputBucket string

	// Optional config
	IncludedTools []string
	ExcludedTools []string
}

func (c Config) validate() error {
	if c.ServiceName == "" {
		return errors.New("ServiceName is required")
	}
	if c.GitHubToken == "" {
		return errors.New("GitHubToken is required")
	}
	if c.ProjectID == "" {
		return errors.New("ProjectID is required")
	}
	if c.Region == "" {
		return errors.New("Region is required")
	}
	if c.SubBuildServiceAccount == "" {
		return errors.New("SubBuildServiceAccount is required")
	}
	if c.SubBuildLogsBucket == "" {
		return errors.New("SubBuildLogsBucket is required")
	}
	if c.SubBuildTestOutputBucket == "" {
		return errors.New("SubBuildTestOutputBucket is required")
	}
	return nil
}
