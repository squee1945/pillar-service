package service

import (
	"net/http"

	"github.com/google/go-github/v62/github"
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
		s.Log.Debug(ctx, "Received Push Event from repo: %s, commit URL: %s, pushed by: %s", *event.Repo.FullName, *event.HeadCommit.URL, *event.Pusher.Name)

	case *github.PullRequestEvent:
		s.Log.Debug(ctx, "Received Pull Request Event from repo: %s, %s action on PR #%d in %s", *event.Repo.FullName, *event.Action, *event.Number, *event.Repo.FullName)

	default:
		s.Log.Info(ctx, "Received unhandled event type: %s", github.WebHookType(r))
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
