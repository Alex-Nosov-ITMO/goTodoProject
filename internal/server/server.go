package server

import (
	"context"
	"net/http"
)

type Server struct {
	http *http.Server
}

// Run http server
func (s *Server) Run(port string, h http.Handler) error{
	s.http = &http.Server{
		Addr:    ":" + port,
		Handler: h,
	}

	return s.http.ListenAndServe()
}

// Shutdown http server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}