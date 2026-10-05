package run

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"worker_pool/internal/app/conf"
	workerpool "worker_pool/internal/app/worker_pool"
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

	cfg := conf.NewConfEnv()

	wp := workerpool.NewWorkerPool(ctx, cfg)
	wp.Run()

	// Waiting termination signals
	<-ctx.Done()
	log.Printf("Received termination signal, the application will be stopped...")

	wp.Stop()
}
