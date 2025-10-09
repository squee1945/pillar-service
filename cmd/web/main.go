package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sethvargo/go-envconfig"
	"github.com/squee1945/pillar-service/pkg/logger"
	"github.com/squee1945/pillar-service/pkg/secrets"
	"github.com/squee1945/pillar-service/pkg/service"
)

type config struct {
	// KmsKeyName              string `env:"KMS_KEY_NAME,required"`
	// RunnerServiceAccount    string `env:"RUNNER_SERVICE_ACCOUNT,required"`
	// RunnerImage             string `env:"RUNNER_IMAGE,required"`
	// PromptBucket            string `env:"PROMPT_BUCKET,required"`
	GitHubAppID                int64  `env:"GITHUB_APP_ID,required"`
	GitHubWebhookSecretName    string `env:"GITHUB_WEBHOOK_SECRET_NAME,required"`
	GitHubPrivateKeySecretName string `env:"GITHUB_PRIVATE_KEY_SECRET_NAME,required"`
	// GeminiApiKeyFile        string `env:"GEMINI_API_KEY_FILE,required"`

	Port           string        `env:"PORT,default=8080"`
	SecretCacheTTL time.Duration `env:"SECRET_CACHE_TTL,default=1m"`
	// RunnerTimeout          time.Duration `env:"RUNNER_TIMEOUT,default=180s"`
	TokenExchangeTimeout time.Duration `env:"TOKEN_EXCHANGE_TIMEOUT,default=30s"`
	// RepoPermissionTimeout  time.Duration `env:"REPO_PERMISSION_TIMEOUT,default=30s"`
	// McpToolTimeout         time.Duration `env:"MCP_TOOL_TIMEOUT,default=30s"`
	// WebhookURLPath         string        `env:"WEBHOOK_URL_PATH,default=/webhook"`
	// DefaultMaxSessionTurns int           `env:"DEFAULT_MAX_SESSION_TURNS,default=20"`
}

func main() {
	ctx := context.Background()
	log := logger.New()

	var c config
	if err := envconfig.Process(ctx, &c); err != nil {
		fail(ctx, log, "processing environment variables: %v", err)
	}

	secretAccessor, err := secrets.New(ctx, c.SecretCacheTTL)
	if err != nil {
		fail(ctx, log, "creating secret accessor: %v", err)
	}
	defer secretAccessor.Close()

	serverConfig := service.Config{
		Log:                     log,
		AppID:                   c.GitHubAppID,
		Secrets:                 secretAccessor,
		WebhookSecretName:       c.GitHubWebhookSecretName,
		AppPrivateKeySecretName: c.GitHubPrivateKeySecretName,
		TokenExchangeTimeout:    c.TokenExchangeTimeout,
	}

	server, err := service.New(ctx, serverConfig)
	if err != nil {
		fail(ctx, log, "creating service: %v", err)
	}

	log.Info(ctx, strings.Repeat("=", 120))
	log.Info(ctx, "Starting server on port "+c.Port)
	if err := http.ListenAndServe(":"+c.Port, server.Handler()); err != nil {
		fail(ctx, log, "server failed: %v", err)
	}
}

func fail(ctx context.Context, log logger.L, format string, args ...any) {
	log.Critical(ctx, "FAILED: "+format, args...)
	os.Exit(1)
}
