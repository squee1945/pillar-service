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

	switch action := event.GetAction(); action {
	case "published":
		break
	default:
		s.Log.Debug(ctx, "Ignoring release %s event", action)
		return nil
	}

	// Find dependents (reverse dependencies)
	dependent := repository{owner: "kmonty-catamaran", repo: "deps-weather-webapp"}

	// Fork them
	newOwner, newRepo, err := s.fork(ctx, event.GetInstallation(), dependent.owner, dependent.repo, event.GetRepo().GetOwner())
	if err != nil {
		return fmt.Errorf("forking %s/%s: %v", dependent.owner, dependent.repo, err)
	}

	s.Log.Info(ctx, "Forked %s/%s to %s/%s", dependent.owner, dependent.repo, newOwner, newRepo)

	// TODO: Launch LLM to update and test code
	// TODO: Create pull request from fork back to dependency

	return nil
}
