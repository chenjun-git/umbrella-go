package handler

import (
	"github.com/go-chi/chi"

	"umbrella-go/handler/v1_0"
)

func RegisterBackendRouter() *chi.Mux {
	router := chi.NewRouter()
	registerRouter(router)
	return router
}

func registerRouter(r chi.Router) {
	v1_0.RegisterRouter(r)
}
