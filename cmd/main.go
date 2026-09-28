package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"worker_pool/internal/core"
	"worker_pool/internal/server"
)

func main() {
	// Graceful Shutdown via NotifyContext
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,    // Ctrl+C
		syscall.SIGTERM, // docker / k8s
		syscall.SIGQUIT, // debug stack dump
	)
	defer stop()

	// Чтение дополнительных переменных окружения
	TaskProbabilityFailedEnvRead()
	TaskWorkImitationEnvRead()

	// Создание контроллера для взаимодействия с kubernetes
	ka := core.NewKubernetesAdapter()

	// Создание менеджера состояний задач
	tsm := core.NewTaskStateManager()

	// Создаем пулл воркеров
	wp := core.NewWorkerPool(QueueSizeEnvRead(), WorkersNumEnvRead(), tsm)
	wp.Run()

	// Запуск сервера
	myServer := server.NewServer()
	myServer.Init(tsm, wp, ka)

	// Сигнал для kubernetes о готовности
	ka.SetStateReadyOK()

	// Создаем канал с сигналом об остановки сервиса
	osSignalsChan := make(chan os.Signal, 1)
	signal.Notify(osSignalsChan, os.Interrupt, syscall.SIGTERM)

	// Ждем сигнал об остановке
	<-ctx.Done()
	log.Printf("Received termination signal, the application will be stopped...")

	// Сигнал для kubernetes о скором отклбчении приложения и остановка воркеров
	ka.SetStateNotReady()
	wp.Stop()
}

func QueueSizeEnvRead() int {
	// Считываем размер очереди
	QUEUE_SIZE := os.Getenv("QUEUE_SIZE")
	if QUEUE_SIZE == "" {
		log.Printf("queue size not set! selected default value %v", 64)
		return 64
	}

	queue_size, err := strconv.Atoi(QUEUE_SIZE)
	if err != nil {
		log.Printf("invalid queue size! selected default value %v", 64)
		return 64
	}

	if queue_size <= 0 {
		log.Printf("invalide queue size! value must be more than zero. selected default value %v", 64)
		return 64
	}

	return queue_size
}

func WorkersNumEnvRead() int {
	// Считываем количество воркеров
	WORKERS := os.Getenv("WORKERS")
	if WORKERS == "" {
		log.Printf("workers not set! selected default value %v", 4)
		return 4
	}

	workers, err := strconv.Atoi(WORKERS)
	if err != nil {
		log.Printf("invalid workers value! selected default value %v", 4)
		return 4
	}

	if workers <= 0 {
		log.Printf("invalide workers value! value must be more than zero. selected default value %v", 4)
		return 4
	}

	return workers
}

func TaskProbabilityFailedEnvRead() {
	val := os.Getenv("TASK_PROBABILITY_FAILED")
	if val == "" {
		return
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return
	}

	if intVal <= 0 {
		return
	}

	core.TASK_PROBABILITY_FAILED = intVal
}

func TaskWorkImitationEnvRead() {
	val := os.Getenv("TASK_WORK_IMITATION_DURATION")
	if val == "" {
		return
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return
	}

	if intVal <= 0 {
		return
	}

	core.TASK_WORK_IMITATION_DURATION = intVal
}
