package service

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/squee1945/pillar-service/pkg/logger"
	"github.com/squee1945/pillar-service/pkg/secrets"
)

const (
	defaultServiceTag = "pillar"
)

type Config struct {
	Log logger.L

	AppID int64

	Secrets                 *secrets.S
	WebhookSecretName       string
	AppPrivateKeySecretName string

	TokenExchangeTimeout time.Duration

	// Optional
	Transport  http.RoundTripper
	ServiceTag string
}

type Service struct {
	Config
}

func New(ctx context.Context, cfg Config) (*Service, error) {
	if cfg.Transport == nil {
		cfg.Transport = http.DefaultTransport
	}
	if cfg.ServiceTag == "" {
		cfg.ServiceTag = defaultServiceTag
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
	w.Write([]byte("Hello"))
}

const consonants = "bcdfghjklmnpqrstvwxyz"

func randomString(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	charsetLength := len(consonants)
	for i := range b {
		b[i] = consonants[rand.Intn(charsetLength)]
	}
	return string(b)
}
