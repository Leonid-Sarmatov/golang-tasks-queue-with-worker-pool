package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	workerpool "worker_pool/internal/app/worker_pool"
)

type KuberMetricsInput struct {
}

type KuberMetricsOutput struct {
}

func MakeHandler_KuberMetrics(wp *workerpool.WorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[KuberMetrics]")

		metrics := map[string]interface{}{
			"worker_pool_queue_length": wp.QueueLoad(),
		}
		json.NewEncoder(w).Encode(metrics)

		w.WriteHeader(http.StatusOK)
	}
}
