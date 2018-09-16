package httpcaller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"umbrella-go/umbrella-common/caller"
	"umbrella-go/umbrella-common/middleware/http"
)

func newEchoServer(middlewares ...httpmiddleware.ServerMiddleware) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content := r.URL.Query().Get("content")

		w.Header().Set("contentType", "text/plain")
		w.Write([]byte(content))
	})

	h := httpmiddleware.WithServerMiddleware(handler, middlewares...)
	return httptest.NewServer(h)
}

func assertCallerName(t *testing.T, name string) httpmiddleware.ServerMiddleware {
	return func(rw http.ResponseWriter, req *http.Request, next http.Handler) {
		callerName := caller.CallerNameFromContext(req.Context())
		assert.Equal(t, name, callerName)
		next.ServeHTTP(rw, req)
	}
}

func TestCallerName(t *testing.T) {
	server := newEchoServer(ExtractCallerName(), assertCallerName(t, "test"))
	defer server.Close()

	client := &http.Client{
		Transport: InjectCallerName("test").Wrap(http.DefaultTransport),
	}

	client.Get(server.URL)
}
