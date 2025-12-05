package service

import (
	"fmt"
	"net/http"

	"github.com/squee1945/pillar-service/pkg/logger"
	"github.com/squee1945/pillar-service/pkg/secrets"
)

type Config struct {
	Log logger.L

	ProjectID string
	Region    string
	AppID     int64

	Secrets                 *secrets.S
	WebhookSecretName       string
	AppPrivateKeySecretName string
	GeminiAPIKeySecretName  string

	SubBuildServiceAccount   string
	SubBuildLogsBucket       string
	SubBuildTestOutputBucket string

	// Optional
	Transport   http.RoundTripper
	ServiceName string
}

func (c Config) validate() error {
	if c.ProjectID == "" {
		return fmt.Errorf("ProjectID must be set")
	}
	if c.Region == "" {
		return fmt.Errorf("Region must be set")
	}
	if c.AppID == 0 {
		return fmt.Errorf("AppID must be set")
	}
	if c.Secrets == nil {
		return fmt.Errorf("Secrets must be set")
	}
	if c.WebhookSecretName == "" {
		return fmt.Errorf("WebhookSecretName must be set")
	}
	if c.AppPrivateKeySecretName == "" {
		return fmt.Errorf("AppPrivateKeySecretName must be set")
	}
	if c.GeminiAPIKeySecretName == "" {
		return fmt.Errorf("GeminiAPIKeySecretName must be set")
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
