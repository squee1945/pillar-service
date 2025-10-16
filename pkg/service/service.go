package service

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultTokenExchangeTimeout = 30 * time.Second
	defaultServiceName          = "pillar"
)

type Service struct {
	Config
}

func New(ctx context.Context, cfg Config) (*Service, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	if cfg.Transport == nil {
		cfg.Transport = http.DefaultTransport
	}
	if cfg.ServiceName == "" {
		cfg.ServiceName = defaultServiceName
	}
	if cfg.TokenExchangeTimeout == 0 {
		cfg.TokenExchangeTimeout = defaultTokenExchangeTimeout
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
