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
	t := &promptIssueCommentCreatedPopulatePR{
		projectID:        s.ProjectID,
		region:           s.Region,
		commit:           commit,
		testOutputBucket: s.SubBuildTestOutputBucket,
		event:            event,
	}

	prompt, err := s.renderPrompt(ctx, t)
	if err != nil {
		return fmt.Errorf("rendering prompt: %v", err)
	}

	opts := []adkConfigOption{
		withDevHelperIncludeToolsADK(
			"create_cloud_build",
			"get_cloud_build",
			"get_cloud_build_logs",
			"fetch_test_output",
			"fetch_provenance",
		),
		withGithubIncludeToolsADK(
			"add_issue_comment",
			"get_commit",
			"get_file_contents",
			// "list_branches",
			"list_commits",
			"pull_request_read",
			"search_code",
			// "update_pull_request",
		),
	}

	if err := s.runADK(ctx, installationID, event.GetRepo(), prompt, opts...); err != nil {
		return fmt.Errorf("running ADK: %v", err)
	}

	return nil
}

func (s *Service) extractServiceCommand(_ context.Context, body string) (string, bool) {
	cmd, forService := strings.CutPrefix(strings.TrimSpace(body), "/"+s.ServiceName+" ")
	cmd = strings.TrimSpace(cmd)
	return cmd, forService
}
