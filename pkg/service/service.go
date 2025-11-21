package service

import (
	"context"
	"fmt"
	"net/http"
	"text/template"
)

const (
	defaultServiceName = "pillar"
)

type Service struct {
	Config

	prompts *template.Template
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

	prompts, err := parsePromptTemplates(ctx)
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	return &Service{Config: cfg, prompts: prompts}, nil
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
