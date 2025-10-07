package service

import (
	"fmt"
	"net/http"
)

func (s *Service) serverError(w http.ResponseWriter, r *http.Request, status int, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	s.Log.Error(r.Context(), format, args...)
	w.WriteHeader(status)
	w.Write([]byte("Server error: " + msg))
}

func (s *Service) clientError(w http.ResponseWriter, r *http.Request, status int, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	s.Log.Debug(r.Context(), format, args...)
	w.WriteHeader(status)
	w.Write([]byte("Client error: " + msg))
}
