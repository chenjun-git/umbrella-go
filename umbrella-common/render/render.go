package render

import (
	"net/http"

	chiRender "github.com/go-chi/render"

	"umbrella-go/umbrella-common/errors"
	"umbrella-go/umbrella-common/lang"
)

type RenderFunc func(w http.ResponseWriter, r *http.Request, v interface{})
type ErrorMsgGetter func(code int, languages []string) string

// 构造一个JSON Render
// 在render之前，先通过错误码获取message信息并填充到Error结构中
func MakeJSON(errorMsgGetter ErrorMsgGetter) RenderFunc {
	return func(w http.ResponseWriter, r *http.Request, v interface{}) {
		languages := lang.FromOutgoingContext(r.Context())

		if err, ok := v.(errors.Error); ok {
			// TODO monitor RequestWithRespCode
			if err.GetMessage() == "" {
				if msg := errorMsgGetter(err.GetCode(), languages); msg != "" {
					err.SetMessage(msg)
				} else {
					err.SetMessage("Unknown error")
				}
			}
		}

		chiRender.JSON(w, r, v)
	}
}
