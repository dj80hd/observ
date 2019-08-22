package app

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics for USE obsrvability
type Metrics struct {
	utilization prometheus.Gauge
	saturation  prometheus.Gauge
	errors      prometheus.Counter
}

// App contains the jobs queue and where to publish metrics
type App struct {
	jobs    chan Job
	metrics Metrics
}

// New is a constructor
func New(workers int) *App {
	a := &App{}

	a.jobs = make(chan Job, workers)

	a.metrics = Metrics{
		promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "observ_worker_utilization",
				Help: "number of busy workers",
			},
		),
		promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "observ_worker_saturation",
				Help: "number of queued jobs",
			},
		),
		promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "observ_worker_errors",
				Help: "failed job count",
			},
		),
	}

	for i := 0; i < workers; i++ {
		a.createWorker()
	}
	return a
}

// handleReq handles /req api endpoint
func (a *App) handleReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method.", 405)
		return
	}

	duration, err := strconv.Atoi(r.URL.Query().Get("duration"))
	if err != nil {
		duration = 0
	}

	httpcode, err := strconv.Atoi(r.URL.Query().Get("httpcode"))
	if err != nil {
		httpcode = 200
	}

	worksecs, err := strconv.Atoi(r.URL.Query().Get("worksecs"))
	if err != nil {
		worksecs = 0
	}

	workfail := "true" == r.URL.Query().Get("workfail")

	if worksecs > 0 && !a.enqueue(Job{worksecs, workfail}) {
		w.WriteHeader(507)
		return
	}

	time.Sleep(time.Duration(duration) * time.Millisecond)
	w.WriteHeader(httpcode)
}

// createWorker creates a worker thread
func (a *App) createWorker() {
	go func() {
		for {
			j := <-a.jobs // blocks for more jobs

			// When job taken off queue, dec saturation
			a.metrics.saturation.Dec()

			// When worker thread is active, inc utilization
			a.metrics.utilization.Inc()

			if ok := j.Run(); !ok {
				// When a job fails, inc errors
				a.metrics.errors.Inc()
			}

			// When worker thread goes inactive, dec utilization
			a.metrics.utilization.Dec()
		}
	}()
}

// enqueue enqueues a Job or returns false if empty or full queue
func (a *App) enqueue(job Job) bool {
	for {
		select {
		case a.jobs <- job:
			// When a job is put on queue, inc saturation
			a.metrics.saturation.Inc()
			return true
		default:
			return false
		}
	}
}

// ServeHTTP makes App an http.Handler
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {

	case strings.Contains(r.URL.Path, "/metrics"):
		promhttp.Handler().ServeHTTP(w, r)

	case strings.Contains(r.URL.Path, "/req"):
		// InstrumentHandler provides RED metrics automatically
		prometheus.InstrumentHandler(
			"/req",
			http.HandlerFunc(a.handleReq))(w, r)

	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
