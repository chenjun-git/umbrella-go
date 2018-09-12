package v1_0

import (
	"github.com/go-chi/chi"
)

func RegisterRouter(r chi.Router) {
	r.Get("/news", NewsHandler)
}