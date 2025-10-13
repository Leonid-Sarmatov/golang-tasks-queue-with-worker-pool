package server

import (
	"log"
	"worker_pool/internal/core"
	"net/http"
	//"fmt"
	"encoding/json"
)

func KubernetesMetricsHandler(tsm *core.TaskStateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HANDLER: kubernetesMetricsHandler")

		metrics := map[string]interface{}{
            "worker_pool_queue_length": tsm.GetTaslListSize(),
        }
        json.NewEncoder(w).Encode(metrics)
	}
}

func KubernetesHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HANDLER: kubernetesHealthHandler")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func KubernetesReadyHandler(ka *core.KubernetesAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HANDLER: kubernetesReadyHandler")

		if ka.IsApplicationReady() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ready"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("shutting_down"))
		}
	}
}

