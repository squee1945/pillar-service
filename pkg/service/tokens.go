package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v75/github"
)

const (
	tokenExchangeTimeout = 30 * time.Second
)

func (s *Service) githubClient(ctx context.Context, installationID int64) (*github.Client, error) {
	privateKey, err := s.Secrets.Read(ctx, s.AppPrivateKeySecretName)
	if err != nil {
		return nil, fmt.Errorf("reading private key: %v", err)
	}

	tr, err := ghinstallation.New(s.Transport, s.AppID, installationID, privateKey)
	if err != nil {
		return nil, fmt.Errorf("creating transport: %v", err)
	}
	return github.NewClient(&http.Client{Transport: tr}), nil
}

func (s *Service) installationToken(ctx context.Context, installationID int64, opts ...installationTokenOption) (string, error) {
	privateKey, err := s.Secrets.Read(ctx, s.AppPrivateKeySecretName)
	if err != nil {
		return "", fmt.Errorf("reading private key: %v", err)
	}

	tr, err := ghinstallation.NewAppsTransport(s.Transport, s.AppID, privateKey)
	if err != nil {
		return "", fmt.Errorf("creating transport: %v", err)
	}

	ghClient := github.NewClient(&http.Client{Transport: tr})

	ctx, cancel := context.WithTimeout(context.Background(), tokenExchangeTimeout)
	defer cancel()

	itops := &github.InstallationTokenOptions{}
	for _, opt := range opts {
		opt(itops)
	}

	token, resp, err := ghClient.Apps.CreateInstallationToken(ctx, installationID, itops)
	defer resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("creating installation token: %v", err)
	}
	if token == nil {
		return "", fmt.Errorf("installation token is nil")
	}
	if token.Token == nil {
		return "", fmt.Errorf("installation token.Token is nil")
	}
	return *token.Token, nil
}

type installationTokenOption func(*github.InstallationTokenOptions)

func withRepoIDs(repoIDs ...int64) installationTokenOption {
	return func(opts *github.InstallationTokenOptions) {
		opts.RepositoryIDs = append(opts.RepositoryIDs, repoIDs...)
	}
}

func withPermissions(permissions *github.InstallationPermissions) installationTokenOption {
	return func(opts *github.InstallationTokenOptions) {
		opts.Permissions = permissions
	}
}
