package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v75/github"
)

const (
	cmdPopulatePR = "populate-pr"
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

	prompt, err := s.renderPrompt(ctx, &promptReleasePublished{event: event, dependent: fmt.Sprintf("https://github.com/%s/%s", dependent.owner, dependent.repo)})
	if err != nil {
		return fmt.Errorf("rendering prompt: %v", err)
	}

	return s.run(ctx, event.GetInstallation().GetID(), fork, prompt)
}

func (s *Service) issueCommentHandler(ctx context.Context, event *github.IssueCommentEvent) error {
	switch action := event.GetAction(); action {
	case "created":
		break
	case "deleted", "edited":
		s.Log.Debug(ctx, "Ignoring issue_comment %s event", action)
		return nil
	}

	installationID := event.GetInstallation().GetID()
	owner, repo := event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName()
	issueID, commentID := event.GetIssue().GetID(), event.GetComment().GetID()

	if !event.GetIssue().IsPullRequest() {
		s.Log.Debug(ctx, "Ignoring issue_comment %d (issue %d, repo %s/%s) created event for non-pull-request.", commentID, issueID, owner, repo)
		return nil
	}

	command, forService := s.extractServiceCommand(ctx, event.GetComment().GetBody())
	if !forService {
		s.Log.Debug(ctx, "Ignoring comment %d (issue %d, repo %s/%s); no service command found.", commentID, issueID, owner, repo)
		return nil
	}

	switch command {
	case cmdPopulatePR:
		break
	default:
		s.Log.Debug(ctx, "Ignoring comment %d (issue %d, repo %s/%s); unknown service command %q.", commentID, issueID, owner, repo, command)
		return nil
	}

	ghClient, err := s.githubClient(ctx, installationID)
	if err != nil {
		return fmt.Errorf("creating github client: %v", err)
	}

	// Update the issue comment emoji to "looking".
	if _, _, err := ghClient.Reactions.CreateIssueCommentReaction(ctx, owner, repo, commentID, "eyes"); err != nil {
		s.Log.Warn(ctx, "Failed to add 'eyes' reaction to comment %d (issue %d, repo %s/%s), continuing: %v", commentID, issueID, owner, repo, err)
	}

	// Fetch the PR head commit.
	issueNum := event.GetIssue().GetNumber()
	pr, _, err := ghClient.PullRequests.Get(ctx, owner, repo, issueNum)
	if err != nil {
		return fmt.Errorf("getting pull request %d: %w", issueNum, err)
	}
	commit := pr.GetHead().GetSHA()

	// Run the prompt.
	prompt, err := s.renderPrompt(ctx, &promptIssueCommentCreatedPopulatePR{commit: commit, event: event})
	if err != nil {
		return fmt.Errorf("rendering prompt: %v", err)
	}
	devHelperIncludeTools := []string{"create_cloud_build", "get_cloud_build", "get_cloud_build_logs"}
	if err := s.run(ctx, installationID, event.GetRepo(), prompt, withDevHelperIncludeTools(devHelperIncludeTools)); err != nil {
		return err
	}

	return nil
}

func (s *Service) extractServiceCommand(_ context.Context, body string) (string, bool) {
	cmd, forService := strings.CutPrefix(strings.TrimSpace(body), "/"+s.ServiceName+" ")
	cmd = strings.TrimSpace(cmd)
	return cmd, forService
}
