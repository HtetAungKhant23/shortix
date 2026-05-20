package api

import (
	"github.com/HtetAungKhant23/shortix/shortener"
	"github.com/go-chi/chi"
)

func GetRootRouter(svc *shortener.Service) *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/api", getApiRouter(svc))

	return r
}

func getApiRouter(svc *shortener.Service) *chi.Mux {
	h := NewHandler(svc)

	api := chi.NewRouter()

	v1 := chi.NewRouter()

	v1.Get("/health", h.HealthCheckHandler)
	v1.Post("/shorten", h.CreateShortURL)

	api.Mount("/v1", v1)

	return api
}
