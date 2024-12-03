package handlers

import (
	"log/slog"
	"main/internal/controller"
	"net/http"
)

type handlers struct {
	logger     *slog.Logger
	controller controller.Controller
}

type Handlers interface {
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)

	Link(mux *http.ServeMux)
}

func New(controller controller.Controller, logger *slog.Logger) Handlers {
	return &handlers{
		logger:     logger,
		controller: controller,
	}
}

func (h *handlers) Link(mux *http.ServeMux) {
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/refresh", h.Refresh)
}

var headers = []string{
	"X-Real-Ip",
	"X-Forwarded-For",
}

// насколько правильно это тут реализовывать, может в сервис перенести
func (h *handlers) readUserIP(r *http.Request) string {
	for i := range headers {
		h := r.Header.Get(headers[i])
		if h != "" {
			return h
		}
	}

	return r.RemoteAddr
}
