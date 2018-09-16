package httpmiddleware

import (
	"net/http"

	"umbrella-go/umbrella-common/lang"
)

// 从HTTP Header中取出languages，将其设置到Context Metadata中
func RequestContextMetadataSetLanguage(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		languages := lang.FromHttpHeader(r.Header)
		ctx := lang.ContextSetLanguages(r.Context(), languages)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
