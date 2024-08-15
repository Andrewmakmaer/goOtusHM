package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Server struct {
	logg   Logger
	server *http.Server
}

type Handler struct {
	logger Logger
	app    Application
}

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
}

type Application interface {
	CreateEvent(context.Context, string, string, string, string, string, string, string) error
	UpdateEvent(context.Context, string, string, string, string, string, string, string) error
	DeleteEvent(string, string) error
	ListEventDay(string) error
	ListEventWeek(string) error
	ListEventsMonth(string) error
}

func NewServer(logger Logger, app Application, host, port string) *Server {
	addr := net.JoinHostPort(host, port)
	handler := &Handler{logger: logger, app: app}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", loggingMiddleware(handler.Hello))

	newServer := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{server: newServer}
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Error(s.server.ListenAndServe().Error())
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.server.Shutdown(ctx)
	return nil
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello from calendar"))
}
