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

// насколько правильно это тут реализовывать, может в сервис перенести
func readUserIP(h *handlers, r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
		h.logger.Debug("ip addr from r.RemoteAddr, but not from headers",
			slog.Any("ip", IPAddress),
		)
	}
	return IPAddress
}
