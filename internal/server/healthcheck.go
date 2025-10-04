package server

import (
	"fmt"
	"log"
	"net/http"
)

func healthCheckerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Printf("HANDLER: healthCheckerHandler")
	fmt.Fprint(w, "I'm alive!")
}
