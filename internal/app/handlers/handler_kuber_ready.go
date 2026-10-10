package handlers

import (
	"fmt"
	"net/http"

	"worker_pool/internal/app/kuber"
)

func MakeHandler_KuberReady(ka *kuber.KubernetesAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[KuberReady]")

		if ka.IsApplicationReady() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ready"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("shutting_down"))
		}
	}
}
