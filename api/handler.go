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
		helper.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	shortURL, err := h.service.Shorten(req.URL)
	if err != nil {
		if errors.Is(err, shortener.ErrInvalidURL) {
			helper.WriteError(w, http.StatusNotFound, errors.New("invalid or missing URL"))
			return
		}
		helper.WriteError(w, http.StatusInternalServerError, errors.New("failed to create short URL"))
		return
	}

	helper.WriteJSON(w, http.StatusCreated, shortURL)
}

func (h *Handler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	code := helper.GetStringParam(r, "code")

	shortURL, err := h.service.Retrieve(code)
	if err != nil {
		if errors.Is(err, shortener.ErrNotFound) {
			helper.WriteError(w, http.StatusNotFound, errors.New("short URL not found"))
			return
		}
		helper.WriteError(w, http.StatusInternalServerError, errors.New("failed to retrieve short URL"))
	}

	helper.WriteJSON(w, http.StatusOK, shortURL)
}

func (h *Handler) UpdateShortURL(w http.ResponseWriter, r *http.Request) {
	code := helper.GetStringParam(r, "code")

	req := &shortener.UpdateRequest{}
	if err := helper.ReadJSON(r, req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, errors.New("invaild request body"))
		return
	}

	updatedShortURL, err := h.service.Update(code, req.URL)
	if err != nil {
		if errors.Is(err, shortener.ErrNotFound) {
			helper.WriteError(w, http.StatusNotFound, errors.New("short URL not found"))
			return
		}
		helper.WriteError(w, http.StatusInternalServerError, errors.New("failed to update with new URL"))
		return
	}

	helper.WriteJSON(w, http.StatusOK, updatedShortURL)
}

func (h *Handler) DeleteShortURL(w http.ResponseWriter, r *http.Request) {
	code := helper.GetStringParam(r, "code")

	if err := h.service.Delete(code); err != nil {
		if errors.Is(err, shortener.ErrNotFound) {
			helper.WriteError(w, http.StatusNotFound, errors.New("short URL not found"))
			return
		}
		helper.WriteError(w, http.StatusInternalServerError, errors.New("failed to delete short URL"))
		return
	}

	helper.WriteJSON(w, http.StatusNoContent, struct{}{})
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().String(),
	})
}
