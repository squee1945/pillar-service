package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/squee1945/pillar-service/pkg/logger"
	"github.com/squee1945/pillar-service/pkg/secrets"
)

type Config struct {
	Log logger.L

	ProjectID string
	Region    string
	AppID     int64

	PromptBucket         string
	KMSKeyName           string
	RunnerServiceAccount string
	PrepImage            string
	PromptImage          string

	Secrets                 *secrets.S
	WebhookSecretName       string
	AppPrivateKeySecretName string
	GeminiAPIKeySecretName  string

	// Optional
	Transport             http.RoundTripper
	ServiceName           string
	TokenExchangeTimeout  time.Duration
	RunnerTimeout         time.Duration
	MCPToolTimeout        time.Duration
	GeminiMaxSessionTurns int
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
	if c.PromptBucket == "" {
		return fmt.Errorf("PromptBucket must be set")
	}
	if c.KMSKeyName == "" {
		return fmt.Errorf("KMSKeyName must be set")
	}
	if c.RunnerServiceAccount == "" {
		return fmt.Errorf("RunnerServiceAccount must be set")
	}
	if c.PrepImage == "" {
		return fmt.Errorf("PrepImage must be set")
	}
	if c.PromptImage == "" {
		return fmt.Errorf("PromptImage must be set")
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
		return fmt.Errorf("AppPrivateKeySecretName must be set")
	}
	if c.TokenExchangeTimeout == 0 {
		return fmt.Errorf("TokenExchangeTimeout must be set")
	}
	return nil
}
