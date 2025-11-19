package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/v75/github"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	maxForkWaitAttempts = 10
	waitForFork         = 5 * time.Second
	serviceName         = "pillar"
)

type prepDevEnvInput struct {
	UpstreamRepoName string `json:"upstream_repo_name" jsonschema:"The name of the upstream repo that you wish to develop against, in the form 'https://github.com/<owner>/<repo>' or simply '<owner>/<repo>'. A fork of this upstream repository will be created."`
}

type prepDevEnvOutput struct {
	UpstreamOwner         string `json:"upstream_owner" jsonschema:"The owner of the upstream repository."`
	UpstreamRepo          string `json:"upstream_repo" jsonschema:"The repo of the upstream repository."`
	ForkOwner             string `json:"fork_owner" jsonschema:"The owner of the fork that was created."`
	ForkRepo              string `json:"fork_repo" jsonschema:"The repo of the fork that was created."`
	DevBranch             string `json:"dev_branch" jsonschema:"The name of the dev branch that was created."`
	UpstreamDefaultBranch string `json:"upstream_default_branch" jsonschema:"The default branch of the upstream repository, useful when creating contribution pull requests."`
}

func prepDevEnvTool(githubToken string) mcp.ToolHandlerFor[prepDevEnvInput, prepDevEnvOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input prepDevEnvInput) (*mcp.CallToolResult, prepDevEnvOutput, error) {
		source := strings.TrimPrefix(input.UpstreamRepoName, "https://github.com/")
		parts := strings.Split(source, "/")
		if len(parts) != 2 {
			return nil, prepDevEnvOutput{}, fmt.Errorf("invalid upstream repo %q, expect <owner>/<repo>", input.UpstreamRepoName)
		}
		upstreamOwner, upstreamRepo := parts[0], parts[1]

		// Create a fork remotely.
		f, err := fork(ctx, githubToken, upstreamOwner, upstreamRepo)
		if err != nil {
			return nil, prepDevEnvOutput{}, err
		}

		// Clone the fork to local enviroment.
		forkOwner, forkRepo := f.GetOwner().GetLogin(), f.GetName()
		if err := clone(ctx, githubToken, forkOwner, forkRepo); err != nil {
			return nil, prepDevEnvOutput{}, err
		}

		// Checkout a dev branch in the local enviroment.
		devBranch, err := checkoutDevBranch(ctx, f.GetName())
		if err != nil {
			return nil, prepDevEnvOutput{}, err
		}

		// Find the default branch of the original repo (for a future pull request).
		defaultBranch, err := defaultBranch(ctx, githubToken, upstreamOwner, upstreamRepo)
		if err != nil {
			return nil, prepDevEnvOutput{}, err
		}

		output := prepDevEnvOutput{
			UpstreamOwner:         upstreamOwner,
			UpstreamRepo:          upstreamRepo,
			ForkOwner:             forkOwner,
			ForkRepo:              forkRepo,
			DevBranch:             devBranch,
			UpstreamDefaultBranch: defaultBranch,
		}

		return nil, output, nil
	}
}

func fork(ctx context.Context, githubToken, owner, repo string) (*github.Repository, error) {
	ghClient := github.NewClient(nil).WithAuthToken(githubToken)

	opts := &github.RepositoryCreateForkOptions{DefaultBranchOnly: true} // TODO: handle organization?
	fork, resp, err := ghClient.Repositories.CreateFork(ctx, owner, repo, opts)
	if err != nil {
		if _, ok := err.(*github.AcceptedError); !ok {
			return nil, fmt.Errorf("creating fork (status code: %d): %v", resp.StatusCode, err)
		}
	}

	attempt := 1
	for {
		if attempt > maxForkWaitAttempts {
			return nil, fmt.Errorf("failed to wait for fork after %d attempts", attempt)
		}
		repoObj, resp, err := ghClient.Repositories.Get(ctx, fork.GetOwner().GetLogin(), fork.GetName())
		if err != nil {
			return nil, fmt.Errorf("getting fork: %v", err)
		}
		if resp.StatusCode == http.StatusOK {
			return repoObj, nil
		}
		if resp.StatusCode == http.StatusNotFound {
			time.Sleep(waitForFork)
			attempt++
			continue
		}
		return nil, fmt.Errorf("getting fork (status code: %d): %v", resp.StatusCode, err)
	}
}

func clone(ctx context.Context, githubToken, owner, repo string) error {
	args := []string{
		"clone",
		"--depth", "1",
		fmt.Sprintf("https://x-access-token:%s@github.com/%s/%s.git", githubToken, owner, repo),
	}
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func checkoutDevBranch(ctx context.Context, repo string) (string, error) {
	devBranch := fmt.Sprintf("%s-%d-%s", serviceName, time.Now().Unix(), randomString(4))
	args := []string{
		"switch",
		"-c",
		devBranch,
	}
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = repo
	return devBranch, cmd.Run()
}

func defaultBranch(ctx context.Context, githubToken, owner, repo string) (string, error) {
	ghClient := github.NewClient(nil).WithAuthToken(githubToken)
	repoObj, _, err := ghClient.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return "", fmt.Errorf("getting repo: %v", err)
	}
	return repoObj.GetDefaultBranch(), nil
}

const consonants = "bcdfghjklmnpqrstvwxyz"

func randomString(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	charsetLength := len(consonants)
	for i := range b {
		b[i] = consonants[rand.Intn(charsetLength)]
	}
	return string(b)
}
