package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"worker_pool/internal/core"
)

type AddTaskRequest struct {
	ID         string `json:"id"`
	Payload    string `json:"payload"`
	MaxRetries int    `json:"max_retries"`
}

func GetAddTaskHandler(wp *core.WorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HANDLER: addTaskHandler")

		// Читаем тело запроса
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Failed to read request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Декодируем JSON в структуру
		var request AddTaskRequest
		if err := json.Unmarshal(body, &request); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Добавляем задачу в очередь на обработку
		err = wp.AddTastToQueue(core.NewTask(request.ID, request.Payload, request.MaxRetries))
		if err != nil {
			log.Printf("Error adding task to queue: %v", err)
			http.Error(w, "Error adding task to queuet", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
