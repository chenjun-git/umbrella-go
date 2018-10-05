package v1_0

import (
	"github.com/go-chi/chi"

	"umbrella-go/umbrella-common/monitor"
)

func RegisterRouter(r chi.Router) {
	r.Get("/news", monitor.HttpHandlerWrapper("newHandler", NewsHandler))
	r.Get("/test", monitor.HttpHandlerWrapper("testHandler", TestHandler))
}
