package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AchilleasB/southpark-messaging/go-api-service/internal/app"
	"github.com/AchilleasB/southpark-messaging/go-api-service/internal/domain"
)

type Handler struct {
	app *app.App
	mux *http.ServeMux
}

func NewHandler(a *app.App) http.Handler {
	h := &Handler{
		app: a,
		mux: http.NewServeMux(),
	}
	h.routes()
	return h.mux
}

func (h *Handler) routes() {
	h.mux.HandleFunc("/messages", h.handleMessages)
	h.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}

func (h *Handler) handleMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg domain.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := msg.Validate(); err != nil {
		http.Error(w, "Invalid message: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.app.PublishMessage(r.Context(), &msg); err != nil {
		http.Error(w, "Failed to publish message: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
