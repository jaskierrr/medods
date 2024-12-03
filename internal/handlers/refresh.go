package handlers

import (
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

	var req models.RefreshRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("Error", slog.Any("Error", err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req.User.IP = h.readUserIP(r)

	responseData, err := h.controller.Refresh(r.Context(), req)
	if err != nil {
		h.logger.Error(err.Error())
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
