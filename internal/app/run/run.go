package run

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"worker_pool/internal/app/conf"
	"worker_pool/internal/app/handlers"
	"worker_pool/internal/app/kuber"

	workerpool "worker_pool/internal/app/worker_pool"
	httpserver "worker_pool/internal/transport/http_server"
)

func Run() {
	// Graceful Shutdown via NotifyContext
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,    // Ctrl+C
		syscall.SIGTERM, // docker / k8s
		syscall.SIGQUIT, // debug stack dump
	)
	defer stop()

	// Upload configurations
	cfg := conf.NewConfEnv()

	// Create kubernetes adapter
	ka := kuber.NewKubernetesAdapter()

	// Init and run worker pool
	wp := workerpool.NewWorkerPool(
		ctx,
		cfg.GetWorkersNumber(),
		cfg.GetTaskQueueSize())
	wp.Run()

	// Create and start HTTP server
	server := httpserver.NewServer(
		cfg.GetHttpPort(),
		cfg.GetHttpReadTimeout(),
		cfg.GetHttpWriteTimeout(),
		cfg.GetHttpShutdownTimeout(),
		handlers.MakeHandler_KuberHealthcheck(),
		handlers.MakeHandler_AddTask(wp),
		handlers.MakeHandler_KuberReady(ka),
		handlers.MakeHandler_KuberMetrics(wp))
	server.Init()

	// Waiting termination signals
	<-ctx.Done()
	log.Printf("Received termination signal, the application will be stopped...")

	server.Stop()

	wp.Stop()
}
