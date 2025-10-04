package server

import (
	"encoding/json"
	"log"
	"net/http"
	"worker_pool/internal/core"
)

type TaskStateResponse struct {
	Tasks []core.Task `json:"tasks"`
}

func GetTaskStateHandler(tsm *core.TaskStateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HANDLER: taskStateHandler")

		resp := TaskStateResponse{
			Tasks: tsm.GetTaskList(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding response to JSON: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
