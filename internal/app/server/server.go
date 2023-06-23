package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Host              string
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

type Server struct {
	httpServer *http.Server
}

func NewServer(config Config, handler http.Handler) *Server {
	addr := fmt.Sprintf(":%d", config.Port)
	log.Println(addr)
	s := &Server{httpServer: &http.Server{
		Addr:              addr,
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
		MaxHeaderBytes:    config.MaxHeaderBytes,
		Handler:           handler,
	}}
	return s
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
