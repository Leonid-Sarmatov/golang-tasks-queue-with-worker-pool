package handlers

import (
	"fmt"
	"net/http"
)

type KuberHealthcheckInput struct {
}

type KuberHealthcheckOutput struct {
}

func MakeHandler_KuberHealthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[KuberHealthcheck]")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "I'm alive!")
	}
}
