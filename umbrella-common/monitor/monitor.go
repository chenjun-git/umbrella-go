package monitor

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"umbrella-go/umbrella-common/caller"
)

var MonitorHandlers = make(map[string]http.Handler)

func RegisterHandlers(r *chi.Mux) {
	if len(MonitorHandlers) == 0 {
		//log.Fatal("cannot start when have no handlers")
	}

	for k, v := range MonitorHandlers {
		if k == "/internal/debug" {
			r.Mount("/internal/debug", v)
		} else {
			r.Handle(k, v)
		}
	}
}

func InitAndListen(listenAddr string) {
	r := chi.NewRouter()
	RegisterHandlers(r)

	go func() {
		//log.Fatal()
		http.ListenAndServe(listenAddr, r)
	}()
}

func MonitorInceptorUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		api := info.FullMethod[strings.LastIndex(info.FullMethod, "/")+1:]
		caller := caller.CallerNameFromContext(ctx)
		if caller == "" {
			caller = "unknown"
		}

		start := time.Now()
		defer func() {
			cost := time.Now().Sub(start)
			code := strconv.Itoa(int(grpc.Code(err)))

			if counter, _ := Monitor.Counter(caller, api, code); counter != nil { // TODO: caller
				counter.Inc()
			}

			if timer, _ := Monitor.Timer(caller, api, code); timer != nil {
				timer.Observe(float64(cost / time.Millisecond))
			}
		}()

		resp, err = handler(ctx, req)
		return
	}
}

func HttpHandlerWrapper(api string, handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			cost := time.Now().Sub(start)
			ctx := r.Context()
			code := strconv.Itoa(respCodeFromContext(ctx))
			caller := caller.CallerNameFromContext(ctx)
			if caller == "" {
				caller = "unknown"
			}

			if counter, _ := Monitor.Counter(caller, api, code); counter != nil { // TODO: caller
				counter.Inc()
			}

			if timer, _ := Monitor.Timer(caller, api, code); timer != nil {
				timer.Observe(float64(cost / time.Millisecond))
			}
		}()

		handler(w, r)
		return
	}
}

type httpResponseCodeKey struct{}

func respCodeFromContext(ctx context.Context) int {
	code, ok := ctx.Value(httpResponseCodeKey{}).(int)
	if !ok {
		return 0
	}
	return code
}

func RequestWithRespCode(r *http.Request, code int) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), httpResponseCodeKey{}, code))
}
