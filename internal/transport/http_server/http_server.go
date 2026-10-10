package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type httpServer struct {
	port            string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
	healthHandler   http.HandlerFunc
	enqueueHandler  http.HandlerFunc
	readyHandler    http.HandlerFunc
	metricsHandler  http.HandlerFunc
	server          http.Server
}

func NewServer(
	port string,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	shutdownTimeout time.Duration,
	healthHandler http.HandlerFunc,
	enqueueHandler http.HandlerFunc,
	readyHandler http.HandlerFunc,
	metricsHandler http.HandlerFunc,
) *httpServer {
	return &httpServer{
		port:            port,
		readTimeout:     readTimeout,
		writeTimeout:    writeTimeout,
		shutdownTimeout: shutdownTimeout,
		healthHandler:   healthHandler,
		enqueueHandler:  enqueueHandler,
		readyHandler:    readyHandler,
		metricsHandler:  metricsHandler,
	}
}

func (s *httpServer) Init() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/enqueue", s.enqueueHandler)
	mux.HandleFunc("/ready", s.readyHandler)
	mux.HandleFunc("/metrics", s.metricsHandler)

	s.server = http.Server{
		Addr:         s.port,
		Handler:      mux,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("http server start failed: %v\n", err)
		}
	}()
	fmt.Printf("http server successful started!\n")
}

func (s *httpServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		fmt.Printf("http server close failed: %v\n", err)
	}
	fmt.Printf("http server is turned off\n")
}
