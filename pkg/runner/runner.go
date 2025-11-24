package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	gcsUploadTimeout = 30 * time.Second

	defaultRunnerTimeout         = 20 * time.Minute
	defaultMCPToolTimeout        = 2 * time.Minute
	defaultGeminiMaxSessionTurns = 200

	devHelperCommand = "devhelpermcp"
)

type R struct {
	Config

	tag string
}

func New(ctx context.Context, cfg Config) (*R, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	if cfg.RunnerTimeout == 0 {
		cfg.RunnerTimeout = defaultRunnerTimeout
	}
	if cfg.DevHelperMCPTimeout == 0 {
		cfg.DevHelperMCPTimeout = defaultMCPToolTimeout
	}
	if cfg.GithubMCPTimeout == 0 {
		cfg.GithubMCPTimeout = defaultMCPToolTimeout
	}
	if cfg.GeminiMaxSessionTurns == 0 {
		cfg.GeminiMaxSessionTurns = defaultGeminiMaxSessionTurns
	}

	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("generating UUID: %v", err)
	}

	return &R{Config: cfg, tag: uid.String()}, nil
}

func (r *R) Run(ctx context.Context) error {
	encryptedGitHubToken, err := kmsEncrypt(ctx, r.KMSKeyName, []byte(r.GitHubToken))
	if err != nil {
		return fmt.Errorf("encrypting GitHub token: %v", err)
	}

	promptGCSPath, err := r.preparePrompt(ctx)
	if err != nil {
		return fmt.Errorf("preparing prompt: %v", err)
	}

	settingsGCSPath, err := r.prepareSettings(ctx)
	if err != nil {
		return fmt.Errorf("preparing settings: %v", err)
	}

	build := cloudbuildpb.Build{
		Steps: []*cloudbuildpb.BuildStep{
			{
				Name: r.PrepImage,
				Env: []string{
					"OWNER=" + r.Owner,
					"REPO=" + r.Repo,
					"DEFAULT_BRANCH=" + r.DefaultBranch,
					"DEV_BRANCH=" + r.DevBranch,
				},
				SecretEnv: []string{
					"GITHUB_TOKEN",
				},
			},
		},
		AvailableSecrets: &cloudbuildpb.Secrets{
			Inline: []*cloudbuildpb.InlineSecret{
				{
					KmsKeyName: r.KMSKeyName,
					EnvMap: map[string][]byte{
						"GITHUB_TOKEN": encryptedGitHubToken,
					},
				},
			},
		},
		ServiceAccount: r.ServiceAccount,
		Timeout:        durationpb.New(r.RunnerTimeout),
		Options: &cloudbuildpb.BuildOptions{
			Logging: cloudbuildpb.BuildOptions_CLOUD_LOGGING_ONLY,
		},
	}

	if r.Prompt != "" {
		step := &cloudbuildpb.BuildStep{
			Name: r.PromptImage,
			Dir:  "/workspace",
			Env: []string{
				"REPO=" + r.Repo,
				"PROMPT_GCS_PATH=" + promptGCSPath,
				"SETTINGS_GCS_PATH=" + settingsGCSPath,
			},
			SecretEnv: []string{
				"GEMINI_API_KEY",
			},
		}
		build.Steps = append(build.Steps, step)

		encryptedGeminiApiKey, err := kmsEncrypt(ctx, r.KMSKeyName, []byte(r.GeminiAPIKey))
		if err != nil {
			return fmt.Errorf("encrypting Gemini API key: %v", err)
		}
		build.AvailableSecrets.Inline[0].EnvMap["GEMINI_API_KEY"] = encryptedGeminiApiKey
	} else {
		r.Log.Warn(ctx, "No prompt specified, skipping prompt step.")
	}

	endpoint := fmt.Sprintf("%s-cloudbuild.googleapis.com:443", r.Region)
	client, err := cloudbuild.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("creating Cloud Build client: %v", err)
	}
	defer client.Close()

	req := &cloudbuildpb.CreateBuildRequest{
		ProjectId: r.ProjectID,
		Build:     &build,
	}

	op, err := client.CreateBuild(ctx, req)
	if err != nil {
		return fmt.Errorf("creating Cloud Build build: %v", err)
	}

	r.Log.Info(ctx, "Runner build created successfully, operation: %s", op.Name())
	return nil
}

func (r *R) prepareSettings(ctx context.Context) (string, error) {
	settings := geminiSettings{
		MCPServers: map[string]mcpServerSettings{
			"devHelper": {
				Description: "High level tools to assist in creating contributions to GitHub repositories",
				Command:     devHelperCommand,
				Args: []string{
					"--github_token=" + r.GitHubToken,
					"--project_id=" + r.ProjectID,
					"--region=" + r.Region,
					"--sub_build_service_account=" + r.SubBuildServiceAccount,
					"--sub_build_logs_bucket=" + r.SubBuildLogsBucket,
					"--sub_build_test_output_bucket=" + r.SubBuildTestOutputBucket,
				},
				Timeout:      r.DevHelperMCPTimeout.Milliseconds(),
				IncludeTools: r.DevHelperIncludeTools,
				ExcludeTools: r.DevHelperExcludeTools,
			},
			"github": {
				Description: "Tools to interact with GitHub repositories",
				HTTPURL:     "https://api.githubcopilot.com/mcp/",
				Headers: map[string]string{
					"Authorization": "Bearer " + r.GitHubToken,
				},
				Timeout:      r.GithubMCPTimeout.Milliseconds(),
				IncludeTools: r.GithubIncludeTools,
				ExcludeTools: r.GithubExcludeTools,
			},
		},
		MaxSessionTurns: r.GeminiMaxSessionTurns,
	}
	settingsBytes, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshalling Gemini settings: %v", err)
	}

	filename := fmt.Sprintf("settings-%s.json", r.tag)
	gcsFilename := fmt.Sprintf("gs://%s/%s", r.PromptBucket, filename)
	return gcsFilename, uploadToGCS(ctx, r.PromptBucket, filename, settingsBytes)
}

func (r *R) preparePrompt(ctx context.Context) (string, error) {
	filename := fmt.Sprintf("prompt-%s.json", r.tag)
	gcsFilename := fmt.Sprintf("gs://%s/%s", r.PromptBucket, filename)
	return gcsFilename, uploadToGCS(ctx, r.PromptBucket, filename, []byte(r.Prompt))
}

type geminiSettings struct {
	MCPServers      map[string]mcpServerSettings `json:"mcpServers"`
	MaxSessionTurns int                          `json:"maxSessionTurns,omitempty"`
}

type mcpServerSettings struct {
	Description string `json:"description,omitempty"`

	HTTPURL string            `json:"httpUrl,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`

	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`

	IncludeTools []string `json:"includeTools,omitempty"`
	ExcludeTools []string `json:"excludeTools,omitempty"`

	Timeout int64 `json:"timeout,omitempty"`
}
