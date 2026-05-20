package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/HtetAungKhant23/shortix/api/helper"
	"github.com/HtetAungKhant23/shortix/shortener"
)

type Handler struct {
	service *shortener.Service
}

func NewHandler(svc *shortener.Service) *Handler {
	return &Handler{
		service: svc,
	}
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	req := &shortener.CreateRequest{}

	if err := helper.ReadJSON(r, req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.Shorten(req.URL)
	if err != nil {
		if errors.Is(err, shortener.ErrInvalidURL) {
			http.Error(w, "invalid or missing URL", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to create short URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	b, err := json.Marshal(shortURL)
	if err != nil {
		panic(err)
	}

	w.Write(b)
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().String(),
	})
}
