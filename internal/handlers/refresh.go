package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"main/internal/models"
	"net/http"
)

func (h *handlers) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req := models.RefreshRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("Error",
			slog.Any("Error", err),
		)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req.User.IP = readUserIP(h, r)

	ctx := context.Background()

	responseData, err := h.controller.Refresh(ctx, req)
	if err != nil {
		h.logger.Error("Error refresh token",
			slog.Any("Error", err),
		)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
