package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v75/github"
)

const (
	waitForFork     = time.Second
	maxForkAttempts = 10
)

func (s *Service) fork(ctx context.Context, owner, repo string, installation *github.Installation) (string, string, error) {
	if installation == nil {
		return "", "", fmt.Errorf("installation is nil")
	}
	account := installation.GetAccount()
	if account == nil {
		return "", "", fmt.Errorf("installation.Account is nil")
	}

	newOwner := account.GetLogin()
	if newOwner == "" {
		return "", "", fmt.Errorf("installation.Account.Login is empty")
	}
	newRepo := s.ServiceTag + "-" + repo + "-" + randomString(6)

	opts := &github.RepositoryCreateForkOptions{
		Name:              newRepo,
		DefaultBranchOnly: true,
	}
	if account.GetType() == "Organization" {
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
		_, resp, err := ghClient.Repositories.CreateFork(ctx, owner, repo, opts)
		if err != nil {
			return "", "", fmt.Errorf("creating fork: %v", err)
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(waitForFork)
		attempt++
	}

	return newOwner, newRepo, nil
}
