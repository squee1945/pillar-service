package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v75/github"
)

const (
	waitForFork     = time.Second
	maxForkAttempts = 10
)

func (s *Service) fork(ctx context.Context, installation *github.Installation, owner, repo string, forker *github.User) (string, string, error) {
	// ijson, err := json.MarshalIndent(installation, "", "  ")
	// if err != nil {
	// 	return "", "", fmt.Errorf("marshalling installation: %v", err)
	// }
	// s.Log.Debug(ctx, "Installation:\n%s", ijson)

	if installation == nil {
		return "", "", fmt.Errorf("installation is nil")
	}
	if forker == nil {
		return "", "", fmt.Errorf("forker is nil")
	}

	newOwner := forker.GetLogin()
	if newOwner == "" {
		return "", "", fmt.Errorf("forker.Login is empty")
	}
	newRepo := s.ServiceTag + "-" + repo + "-" + randomString(6)

	opts := &github.RepositoryCreateForkOptions{
		Name:              newRepo,
		DefaultBranchOnly: true,
	}
	if forker.GetType() == "Organization" {
		opts.Organization = newOwner
	}

	ghClient, err := s.githubClient(ctx, *installation.ID)
	if err != nil {
		return "", "", fmt.Errorf("creating github client: %v", err)
	}

	attempt := 1
	for {
		if attempt > maxForkAttempts {
			return "", "", fmt.Errorf("failed to create fork after %d attempts", attempt)
		}
		s.Log.Debug(ctx, "Creating fork %s/%s (attempt %d)", newOwner, newRepo, attempt)
		fork, resp, err := ghClient.Repositories.CreateFork(ctx, owner, repo, opts)
		if err != nil {
			if _, ok := err.(*github.AcceptedError); !ok {
				return "", "", fmt.Errorf("creating fork (status code: %d): %v", resp.StatusCode, err)
			}
			forkJSON, err := json.MarshalIndent(fork, "", "  ")
			if err != nil {
				return "", "", fmt.Errorf("marshalling fork: %v", err)
			}
			s.Log.Debug(ctx, "Fork accepted:\n%s", forkJSON)
			s.Log.Debug(ctx, "Fork response headers: %s", resp.Header)
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(waitForFork)
		attempt++
	}

	return newOwner, newRepo, nil
}
