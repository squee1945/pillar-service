package main

import (
	"context"
	"net/http"
	"os"

	"github.com/sethvargo/go-envconfig"
	"github.com/squee1945/pillar-service/pkg/logger"
	"github.com/squee1945/pillar-service/pkg/service"
)

type config struct {
	// KmsKeyName              string `env:"KMS_KEY_NAME,required"`
	// RunnerServiceAccount    string `env:"RUNNER_SERVICE_ACCOUNT,required"`
	// RunnerImage             string `env:"RUNNER_IMAGE,required"`
	// PromptBucket            string `env:"PROMPT_BUCKET,required"`
	// GitHubAppID             int64  `env:"GITHUB_APP_ID,required"`
	// GitHubWebhookSecretFile string `env:"GITHUB_WEBHOOK_SECRET_FILE,required"`
	// GitHubPrivateKeyFile    string `env:"GITHUB_PRIVATE_KEY_FILE,required"`
	// GeminiApiKeyFile        string `env:"GEMINI_API_KEY_FILE,required"`

	Port string `env:"PORT,default=8080"`
	// SecretCacheTTL         time.Duration `env:"SECRET_CACHE_TTL,default=1m"`
	// RunnerTimeout          time.Duration `env:"RUNNER_TIMEOUT,default=180s"`
	// TokenExchangeTimeout   time.Duration `env:"TOKEN_EXCHANGE_TIMEOUT,default=30s"`
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

	serverConfig := service.Config{
		Log: log,
	}

	server, err := service.New(ctx, serverConfig)
	if err != nil {
		fail(ctx, log, "creating service: %v", err)
	}

	log.Info(ctx, "Starting server on port "+c.Port)
	if err := http.ListenAndServe(":"+c.Port, server.Handler()); err != nil {
		fail(ctx, log, "server failed: %v", err)
	}
}

func fail(ctx context.Context, log logger.L, format string, args ...any) {
	log.Critical(ctx, "FAILED: "+format, args...)
	os.Exit(1)
}
