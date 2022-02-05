package internalhttp

import (
	"context"
	"net/http"
)

type Server struct { // TODO
	srv *http.Server
}

type Logger interface {
	Infof(fmt string, v ...interface{})
}

type Application interface { // TODO
}

func NewServer(addr string, logger Logger, app Application) *Server {
	mux := http.NewServeMux()

	mux.Handle("/hello", loggingMiddleware(http.HandlerFunc(helloHandler), logger))

	return &Server{srv: &http.Server{Addr: addr, Handler: mux}}
}

func (s *Server) Start(ctx context.Context) error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
