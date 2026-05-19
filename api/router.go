package api

import (
	"github.com/go-chi/chi"
)

func GetRootRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/api", getApiRouter())

	return r
}

func getApiRouter() *chi.Mux {
	api := chi.NewRouter()

	v1 := chi.NewRouter()

	api.Mount("/v1", v1)

	return api
}
