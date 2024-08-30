package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"
)

type eventFields struct {
	// ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	StartTime   string `json:"starttime,omitempty"`
	EndTime     string `json:"endtime,omitempty"`
	// UserID       string `json:"userid"`
	CallDuration string `json:"callduration,omitempty"`
	// SearchTime   string `json:"searchby,omitempty"`
}

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
	ListEventDay(string) (string, error)
	ListEventWeek(string) (string, error)
	ListEventsMonth(string) (string, error)
}

func NewServer(logger Logger, app Application, host, port string) *Server {
	addr := net.JoinHostPort(host, port)
	handler := &Handler{logger: logger, app: app}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", loggingMiddleware(handler.Hello))
	mux.HandleFunc("POST /events/{user_id}/{id}", loggingMiddleware(handler.CreateEvent))
	mux.HandleFunc("PUT /events/{user_id}/{id}", loggingMiddleware(handler.UpdateEvent))
	mux.HandleFunc("DELETE /events/{user_id}/{id}", loggingMiddleware(handler.DeleteEvent))
	mux.HandleFunc("GET /events/{user_id}/{by_time}", loggingMiddleware(handler.ListEvent))

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
	s.logg.Info("Server shutdown")
	return nil
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello from calendar"))
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cansel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cansel()

	eventFields := unmarshBody(w, r)

	h.logger.Info("message", "create event request", "userID", r.PathValue("user_id"),
		"eventID", r.PathValue("id"), "title", eventFields.Title,
		"callDuration", eventFields.CallDuration)

	err := h.app.CreateEvent(ctx, r.PathValue("id"), r.PathValue("user_id"), eventFields.Title,
		eventFields.Description, eventFields.StartTime, eventFields.EndTime, eventFields.CallDuration)
	if err != nil {
		h.logger.Error("message", "fail to create event", "userID", r.PathValue("user_id"),
			"eventID", r.PathValue("id"), "title", eventFields.Title)
		errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}
	errorResponse(w, "Success create event", http.StatusOK)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cansel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cansel()

	eventFields := unmarshBody(w, r)

	err := h.app.UpdateEvent(ctx, r.PathValue("id"), r.PathValue("user_id"), eventFields.Title,
		eventFields.Description, eventFields.StartTime, eventFields.EndTime, eventFields.CallDuration)
	if err != nil {
		errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}

	errorResponse(w, "Succes update event", http.StatusOK)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	_, cansel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cansel()

	err := h.app.DeleteEvent(r.PathValue("user_id"), r.PathValue("id"))
	if err != nil {
		errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}
	errorResponse(w, "Success delete event", http.StatusOK)
}

func (h *Handler) ListEvent(w http.ResponseWriter, r *http.Request) {
	_, cansel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cansel()

	var err error
	var result string

	switch r.PathValue("by_time") {
	case "day":
		result, err = h.app.ListEventDay(r.PathValue("user_id"))
	case "week":
		result, err = h.app.ListEventWeek(r.PathValue("user_id"))
	case "month":
		result, err = h.app.ListEventsMonth(r.PathValue("user_id"))
	}

	if err != nil {
		errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}
	errorResponse(w, result, http.StatusOK)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func unmarshBody(w http.ResponseWriter, r *http.Request) eventFields {
	var eventFields eventFields
	headerContType := r.Header.Get("Content-Type")
	if headerContType != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return eventFields
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&eventFields)
	var unmarshalErr *json.UnmarshalTypeError

	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return eventFields
	}

	return eventFields
}
