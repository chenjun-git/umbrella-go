package httpmiddleware

import (
	"net/http"
)

// 注意受限于http.RoundTripper接口，ClientMiddleware必须是线程安全的
// 并且Request除了读和关闭Body之外不允许修改
type ServerMiddleware func(rw http.ResponseWriter, req *http.Request, next http.Handler)

func (m ServerMiddleware) Wrap(next http.Handler) http.Handler {
	if m == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m(w, r, next)
	})
}

// Chain 组合ServerMiddleware m和m2，m在m2之前执行
func (m ServerMiddleware) Chain(m2 ServerMiddleware) ServerMiddleware {
	if m == nil {
		return m2
	}
	if m2 == nil {
		return m
	}

	return func(rw http.ResponseWriter, req *http.Request, next http.Handler) {
		m(rw, req, m2.Wrap(next))
	}
}

// middlewares从左至右调用
func WithServerMiddleware(h http.Handler, middlewares ...ServerMiddleware) http.Handler {
	if len(middlewares) == 0 {
		return h
	}

	result := h
	for i := len(middlewares) - 1; i >= 0; i-- {
		result = middlewares[i].Wrap(result)
	}
	return result
}

func ChainServerMiddlewares(middlewares ...ServerMiddleware) ServerMiddleware {
	switch len(middlewares) {
	case 0:
		return nil
	case 1:
		return middlewares[0]
	default:
		m := middlewares[len(middlewares)-1]
		for i := len(middlewares) - 2; i >= 0; i-- {
			m = middlewares[i].Chain(m)
		}
		return m
	}
}
