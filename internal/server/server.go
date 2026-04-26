package server

import (
	"log/slog"
	"net/http"
	"time"
)

type Config struct {
	Addr string
}

type Server struct {
	cfg Config
	logger *slog.Logger
	mux  *http.ServeMux
	
}

func NewServer(cfg Config, logger *slog.Logger) *Server {
	mux := http.NewServeMux()
	
	s := &Server{
		cfg: cfg,
		logger: logger,
		mux: mux,
	}
	
s.routes()
return s
	
}


func (s *Server) routes(){
	s.mux.HandleFunc("/healthz", HeathHandler)
}

func (s *Server) HttpServer() * http.Server{
	return &http.Server{
		Addr: s.cfg.Addr,
		Handler: s.mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
