package http

import (
	"context"
	"net/http"
)

type Server struct {
	*http.Server
	name string
}

func NewServer(hs *http.Server, name string) *Server {
	return &Server{
		Server: hs,
		name:   name,
	}
}

func (s *Server) Run() error {
	return s.ListenAndServe()
}

func (s *Server) Name() string {
	return s.name
}

func (s *Server) Close(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
