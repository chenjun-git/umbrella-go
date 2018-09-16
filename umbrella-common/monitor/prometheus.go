package monitor

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var defaultHistogramBuckets = []float64{1, 5, 10, 20, 50, 100, 200, 500, 1000, 2000, 5000, 10000}

var labels = []string{"caller", "api", "code"}

type monitor struct {
	nameSpace string
	subSystem string

	counter *prometheus.CounterVec
	timer   *prometheus.HistogramVec
}

var Monitor *monitor

func init() {
	MonitorHandlers["/internal/metrics"] = promhttp.Handler()
	MonitorHandlers["/internal/ping"] = Monitor.PingHandler()
}

func Init(nameSpace, subSystem string) {
	m := &monitor{
		nameSpace: nameSpace,
		subSystem: subSystem,
	}

	// TODO  分析prometheus原理
	// register api counter
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "count",
			Help:      fmt.Sprintf("api counter for %s system in %s", m.subSystem, m.nameSpace),
		},
		labels,
	)
	prometheus.MustRegister(counter)
	m.counter = counter

	// register api timer
	timer := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: m.nameSpace,
			Subsystem: m.subSystem,
			Name:      "timer",
			Help:      fmt.Sprintf("api timer for %s system in %s", m.subSystem, m.nameSpace),
			Buckets:   defaultHistogramBuckets,
		},
		labels,
	)
	prometheus.MustRegister(timer)
	m.timer = timer

	Monitor = m
}

func (m *monitor) Counter(caller, api, code string) (prometheus.Counter, error) {
	if m.counter == nil {
		return nil, errors.New("no counter registered")
	}

	return m.counter.GetMetricWithLabelValues(caller, api, code)
}

func (m *monitor) Timer(caller, api, code string) (prometheus.Observer, error) {
	if m.timer == nil {
		return nil, errors.New("no timer registered")
	}

	return m.timer.GetMetricWithLabelValues(caller, api, code)
}

// SetVersion 兼容旧的设置版本的方法
func (m *monitor) SetVersion(v Version) {
	InitVersion(v)
}

// GetVersionHandler 兼容旧的GetVersionHandler
func (m *monitor) GetVersionHandler() http.Handler {
	return http.HandlerFunc(GetVersionHandler)
}

func (m *monitor) PingHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		resp := []byte("ok\n")
		w.Write(resp)
	}

	return http.HandlerFunc(fn)
}
