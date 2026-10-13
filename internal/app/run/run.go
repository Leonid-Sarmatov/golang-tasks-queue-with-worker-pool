package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"worker_pool/internal/app/conf"
	"worker_pool/internal/app/execute"
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

	// Create task executor
	execute_func := execute.MakeExecutor(
		cfg.GetTaskProcessingDuration(),
		cfg.GetTaskProbabilityFailed())

	// Create kubernetes adapter
	ka := kuber.NewKubernetesAdapter()

	// Init and run worker pool
	wp := workerpool.NewWorkerPool(
		ctx,
		cfg.GetWorkersNumber(),
		cfg.GetTaskQueueSize(),
		execute_func)
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

	// Mark the application as ready
	ka.SetStateReadyOK()

	// Wait for termination signals
	<-ctx.Done()
	fmt.Printf("Received termination signal, the application will be stopped...\n")

	// Mark the application as not ready
	ka.SetStateNotReady()

	// Stop HTTP server
	server.Stop()

	// Wait for workers and goroutines to exit
	wp.Stop()
}
