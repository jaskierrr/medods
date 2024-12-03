package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"main/internal/models"
	"net/http"
)

func (h *handlers) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := models.User{
		ID: r.URL.Query().Get("id"),
		// насколько правильно это тут реализовывать, может в сервис перенести
		IP: h.readUserIP(r),
	}

	if user.ID == "" {
		h.logger.Error("User ID can't be empty")
		http.Error(w, "User ID can't be empty", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	responseData, err := h.controller.Login(ctx, user)
	if err != nil {
		h.logger.Error("Error",
			slog.Any("Error", err),
		)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		h.logger.Error("Failed to write response",
			slog.Any("Error", err),
		)
		if w.Header().Get("Content-Type") == "" {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}
}
