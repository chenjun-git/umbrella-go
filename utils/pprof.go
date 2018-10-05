package utils

import (
	"net/http"
	"net/http/pprof"
	"runtime"

	"github.com/go-chi/chi"
)

const pprofPrefix = "/debug/pprof"

// PProfHandlers ...
func PProfHandlers() map[string]http.Handler {
	// set only when there's no existing setting
	if runtime.SetMutexProfileFraction(-1) == 0 {
		// 1 out of 5 mutex events are reported, on average
		runtime.SetMutexProfileFraction(5)
	}

	m := make(map[string]http.Handler)

	m[pprofPrefix+"/"] = http.HandlerFunc(pprof.Index)
	m[pprofPrefix+"/profile"] = http.HandlerFunc(pprof.Profile)
	m[pprofPrefix+"/symbol"] = http.HandlerFunc(pprof.Symbol)
	m[pprofPrefix+"/cmdline"] = http.HandlerFunc(pprof.Cmdline)
	m[pprofPrefix+"/trace "] = http.HandlerFunc(pprof.Trace)
	m[pprofPrefix+"/heap"] = pprof.Handler("heap")
	m[pprofPrefix+"/goroutine"] = pprof.Handler("goroutine")
	m[pprofPrefix+"/threadcreate"] = pprof.Handler("threadcreate")
	m[pprofPrefix+"/block"] = pprof.Handler("block")
	m[pprofPrefix+"/mutex"] = pprof.Handler("mutex")

	return m
}

// RegisterPProf ...
func RegisterPProf(r *chi.Mux) {
	pprofHandlers := PProfHandlers()
	if len(pprofHandlers) == 0 {
		panic("RegisterPProf failed!")
	}

	for k, v := range pprofHandlers {
		if k == "/internal/debug" {
			r.Mount("/internal/debug", v)
		} else {
			r.Handle(k, v)
		}
	}
}
