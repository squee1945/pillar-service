package service

import (
	"net/http"

	"github.com/google/go-github/v75/github"
)

func (s *Service) webhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		s.clientError(w, r, http.StatusMethodNotAllowed, "Method %s not allowed", r.Method)
		return
	}

	webhookSecret, err := s.Secrets.Read(ctx, s.WebhookSecretName)
	if err != nil {
		s.serverError(w, r, http.StatusInternalServerError, "reading webhook secret: %v", err)
		return
	}

	payload, err := github.ValidatePayload(r, webhookSecret)
	if err != nil {
		s.clientError(w, r, http.StatusBadRequest, "invalid signature: %v", err)
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		s.clientError(w, r, http.StatusBadRequest, "could not parse webhook: %v", err)
		return
	}

	switch event := event.(type) {

	case *github.PushEvent:
		s.Log.Debug(ctx, "Received push %s event (repo: %q commitURL: %s)", event.GetAction(), event.GetRepo().GetFullName(), event.GetHeadCommit().GetURL())

	case *github.PullRequestEvent:
		s.Log.Debug(ctx, "Received pullRequest %s event (repo: %q pullRequest: %d)", event.GetAction(), event.GetRepo().GetFullName(), event.GetPullRequest().GetNumber())

	case *github.ReleaseEvent:
		s.Log.Debug(ctx, "Received release %s event (repo: %q release: %q)", event.GetAction(), event.GetRepo().GetFullName(), event.GetRelease().GetName())
		if err := s.releaseEventHandler(ctx, event); err != nil {
			s.serverError(w, r, http.StatusInternalServerError, "release event handler: %v", err)
			return
		}

	default:
		s.Log.Info(ctx, "Received unhandled event type: %s", github.WebHookType(r))
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
