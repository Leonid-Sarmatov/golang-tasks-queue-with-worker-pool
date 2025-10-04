package server

import (
	"log"
	"net/http"
	"time"
	"worker_pool/internal/core"
)

type httpServer struct {
	server http.Server
}

func NewServer() *httpServer {
	return &httpServer{}
}

func (s *httpServer) Init(tsm *core.TaskStateManager, wp *core.WorkerPool) {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthCheckerHandler)
	mux.HandleFunc("/status", GetTaskStateHandler(tsm))
	mux.HandleFunc("/enqueue", GetAddTaskHandler(wp))

	s.server = http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Fatalf("<server.go> http server start failed: %v", err)
		}
	}()
	log.Printf("<server.go> http server successful started!")
}
