package httpcaller

import (
	"net/http"

	"umbrella-go/umbrella-common/caller"
	"umbrella-go/umbrella-common/middleware/http"
)

const (
	callerName = "Caller-Name"
)

func InjectCallerName(name string) httpmiddleware.ClientMiddleware {
	return func(req *http.Request, next http.RoundTripper) (*http.Response, error) {
		newReq := addCallerName(req, name)
		return next.RoundTrip(newReq)
	}
}

func addCallerName(req *http.Request, name string) *http.Request {
	newReq := new(http.Request)
	*newReq = *req
	newReq.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		newReq.Header[k] = s
	}
	newReq.Header.Set(callerName, name)
	return newReq
}

func ExtractCallerName() httpmiddleware.ServerMiddleware {
	return func(rw http.ResponseWriter, req *http.Request, next http.Handler) {
		name := req.Header.Get(callerName)
		next.ServeHTTP(rw, req.WithContext(caller.ContextWithCallerName(req.Context(), name)))
	}
}
