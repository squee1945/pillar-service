package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v75/github"
)

const (
	waitForFork         = time.Second
	maxForkWaitAttempts = 10
)

func (s *Service) fork(ctx context.Context, installationID int64, owner, repo string, forker *github.User) (*github.Repository, error) {
	if forker == nil {
		return nil, fmt.Errorf("forker is nil")
	}

	forkOwner := forker.GetLogin()
	if forkOwner == "" {
		return nil, fmt.Errorf("forker.Login is empty")
	}

	opts := &github.RepositoryCreateForkOptions{DefaultBranchOnly: true}
	if forker.GetType() == "Organization" {
		opts.Organization = forkOwner
	}

	ghClient, err := s.githubClient(ctx, installationID)
	if err != nil {
		return nil, fmt.Errorf("creating github client: %v", err)
	}

	s.Log.Debug(ctx, "Creating fork of %s/%s for %s", owner, repo, forkOwner)
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
			s.Log.Debug(ctx, "Fork %s/%s found", fork.GetOwner().GetLogin(), fork.GetName())
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
