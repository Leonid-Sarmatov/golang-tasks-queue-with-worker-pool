package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Task struct {
	ID             string `json:"id"`
	Payload        string `json:"payload"`
	MaxRetries     int    `json:"max_retries"`
	CurrentRetries int    `json:"current_retries"`
	State          string `json:"state"`
}

type TaskStateResponse struct {
	Tasks []Task `json:"tasks"`
}

func main() {
	port := "8080"
	url := fmt.Sprintf("http://localhost:%s/status", port)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var taskResponse TaskStateResponse
	if err := json.Unmarshal(body, &taskResponse); err != nil {
		panic(fmt.Sprintf("Error decoding JSON: %v", err))
	}

	for i, task := range taskResponse.Tasks {
		fmt.Printf("Task %3d:\t", i+1)
		fmt.Printf("ID: %s\t\t", task.ID)
		fmt.Printf("MaxRetries: %3d\t\t", task.MaxRetries)
		fmt.Printf("CurrentRetries: %3d\t", task.CurrentRetries)
		fmt.Printf("State: %s\n", task.State)
	}
}
