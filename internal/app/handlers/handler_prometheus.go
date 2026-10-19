package handlers

import (
	"fmt"
	"log"
	"net/http"

	workerpool "worker_pool/internal/app/worker_pool"
)

func MakeHandler_PrometheusMetrics(wp *workerpool.WorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Content-Type",
			"text/plain; version=0.0.4; charset=utf-8",
		)

		_, err := fmt.Fprintf(w,
			"# HELP worker_pool_queue_length Number of queued tasks.\n"+
				"# TYPE worker_pool_queue_length gauge\n"+
				"worker_pool_queue_length %d\n",
			wp.QueueLoad(),
		)
		if err != nil {
			log.Printf("failed to write Prometheus metrics: %v", err)
		}
	}
}
