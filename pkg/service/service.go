package service

import (
	"context"
	"net/http"
	"time"

	"github.com/squee1945/pillar-service/pkg/logger"
	"github.com/squee1945/pillar-service/pkg/secrets"
)

type Config struct {
	Log logger.L

	AppID int64

	Secrets                 *secrets.S
	WebhookSecretName       string
	AppPrivateKeySecretName string

	TokenExchangeTimeout time.Duration

	// Optional
	Transport http.RoundTripper
}

type Service struct {
	Config
}

func New(ctx context.Context, cfg Config) (*Service, error) {
	if cfg.Transport == nil {
		cfg.Transport = http.DefaultTransport
	}
	return &Service{Config: cfg}, nil
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/webhook", http.HandlerFunc(s.webhook))
	mux.Handle("/", http.HandlerFunc(s.indexHandler))
	return mux
}

func (s *Service) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
}
