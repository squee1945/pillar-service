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
	ProjectID string `env:"PROJECT_ID,required"`
	Region    string `env:"REGION,required"`

	KMSKeyName                 string `env:"KMS_KEY_NAME,required"`
	RunnerServiceAccount       string `env:"RUNNER_SERVICE_ACCOUNT,required"`
	PrepImage                  string `env:"PREP_IMAGE,required"`
	PromptImage                string `env:"PROMPT_IMAGE,required"`
	PromptBucket               string `env:"PROMPT_BUCKET,required"`
	GitHubAppID                int64  `env:"GITHUB_APP_ID,required"`
	GitHubWebhookSecretName    string `env:"GITHUB_WEBHOOK_SECRET_NAME,required"`
	GitHubPrivateKeySecretName string `env:"GITHUB_PRIVATE_KEY_SECRET_NAME,required"`
	GeminiApiKeySecretName     string `env:"GEMINI_API_KEY_SECRET_NAME,required"`

	Port           string        `env:"PORT,default=8080"`
	SecretCacheTTL time.Duration `env:"SECRET_CACHE_TTL,default=1m"`

	// A "SubBuild" is a build that is configured and created by the runner.
	SubBuildServiceAccount   string `env:"SUB_BUILD_SERVICE_ACCOUNT",required`
	SubBuildLogsBucket       string `env:"SUB_BUILD_LOGS_BUCKET",required`
	SubBuildTestOutputBucket string `env:"SUB_BUILD_TEST_OUTPUT_BUCKET",required`
	SubBuildGoRepository     string `env:"SUB_BUILD_GO_REPOSITORY",required`
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
		Log:                      log,
		AppID:                    c.GitHubAppID,
		Secrets:                  secretAccessor,
		WebhookSecretName:        c.GitHubWebhookSecretName,
		AppPrivateKeySecretName:  c.GitHubPrivateKeySecretName,
		ProjectID:                c.ProjectID,
		Region:                   c.Region,
		PromptBucket:             c.PromptBucket,
		KMSKeyName:               c.KMSKeyName,
		RunnerServiceAccount:     c.RunnerServiceAccount,
		PrepImage:                c.PrepImage,
		PromptImage:              c.PromptImage,
		GeminiAPIKeySecretName:   c.GeminiApiKeySecretName,
		SubBuildServiceAccount:   c.SubBuildServiceAccount,
		SubBuildLogsBucket:       c.SubBuildLogsBucket,
		SubBuildTestOutputBucket: c.SubBuildTestOutputBucket,
		SubBuildGoRepository:     c.SubBuildGoRepository,
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
