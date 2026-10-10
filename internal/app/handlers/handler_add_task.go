package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	workerpool "worker_pool/internal/app/worker_pool"
	"worker_pool/internal/core/domain"
)

const maxBodySize int64 = 64 * 1024

type AddTaskRequest struct {
	ID          string `json:"id"`
	Payload     string `json:"payload"`
	MaxAttempts int32  `json:"max_retries"`
}

func MakeHandler_AddTask(wp *workerpool.WorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request", http.StatusBadRequest)
			return
		}

		// Unmarshal request
		var request AddTaskRequest
		if err := json.Unmarshal(body, &request); err != nil {
			fmt.Printf("failed to decoding JSON: %v\n", err)
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Create task
		task, err := domain.NewTask(
			domain.ID(request.ID),
			request.Payload,
			domain.Attempt(request.MaxAttempts))
		if err != nil {
			fmt.Printf("failed to create task: %v\n", err)
			http.Error(w, "failed to create task", http.StatusInternalServerError)
			return
		}

		// Put task into queue
		err = wp.AddTaskToQueue(task)
		if err != nil {
			fmt.Printf("failed to add task to queue: %v\n", err)
			http.Error(w, "failed to add task to queue", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
