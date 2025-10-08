package service

import (
	"context"
	"fmt"

	"github.com/google/go-github/v75/github"
)

type repository struct {
	owner string
	repo  string
}

func (s *Service) releaseEventHandler(ctx context.Context, event *github.ReleaseEvent) error {
	// Find dependents (reverse dependencies)
	dependent := repository{owner: "kmonty-catamaran", repo: "deps-weather-webapp"}

	// Fork them
	newOwner, newRepo, err := s.fork(ctx, dependent.owner, dependent.repo, event.GetInstallation())
	if err != nil {
		return fmt.Errorf("forking %s/%s: %v", dependent.owner, dependent.repo, err)
	}

	s.Log.Info(ctx, "Forked %s/%s to %s/%s", dependent.owner, dependent.repo, newOwner, newRepo)

	// TODO: Launch LLM to update and test code
	// TODO: Create pull request from fork back to dependency

	return nil
}
