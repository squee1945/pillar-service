package service

import (
	"context"
	"fmt"

	"github.com/google/go-github/v75/github"
)

var (
	dependent = repository{owner: "kmonty-catamaran", repo: "deps-weather-webapp"}
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

	// Find dependents (reverse dependencies).
	// TODO

	fork, err := s.fork(ctx, event.GetInstallation().GetID(), dependent.owner, dependent.repo, event.GetRepo().GetOwner())
	if err != nil {
		return fmt.Errorf("forking %s/%s: %v", dependent.owner, dependent.repo, err)
	}

	prompt := `
Using the local checked out repository and the github tools, list the releases for this repository.
`

	return s.run(ctx, event.GetInstallation().GetID(), fork, prompt)
}
