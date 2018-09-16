package httpmiddleware

import (
	"net/http"
)

// 注意受限于http.RoundTripper接口，ClientMiddleware必须是线程安全的
// 并且Request除了读和关闭Body之外不允许修改
type ClientMiddleware func(req *http.Request, next http.RoundTripper) (*http.Response, error)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func (m ClientMiddleware) Wrap(next http.RoundTripper) http.RoundTripper {
	if m == nil {
		return next
	}
	return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return m(req, next)
	})
}

// Chain 组合ClientMiddleware m和m2，m在m2之前执行
func (m ClientMiddleware) Chain(m2 ClientMiddleware) ClientMiddleware {
	if m == nil {
		return m2
	}
	if m2 == nil {
		return m
	}

	return func(req *http.Request, next http.RoundTripper) (*http.Response, error) {
		return m(req, m2.Wrap(next))
	}
}

// middlewares从左至右调用
func WithClientMiddleware(rt http.RoundTripper, middlewares ...ClientMiddleware) http.RoundTripper {
	if len(middlewares) == 0 {
		return rt
	}

	result := rt
	for i := len(middlewares) - 1; i >= 0; i-- {
		result = middlewares[i].Wrap(result)
	}
	return result
}

func ChainClientMiddlewares(middlewares ...ClientMiddleware) ClientMiddleware {
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
