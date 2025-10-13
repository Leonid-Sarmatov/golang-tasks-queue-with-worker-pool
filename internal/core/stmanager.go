package core

import (
	"log"
	"sync"
)

type TaskStateManager struct {
	mu        sync.Mutex
	tasksList map[string]*Task
}

func NewTaskStateManager() *TaskStateManager {
	return &TaskStateManager{
		mu:        sync.Mutex{},
		tasksList: make(map[string]*Task),
	}
}

func (tm *TaskStateManager) GetTaslListSize() int {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	return len(tm.tasksList)
}

func (tm *TaskStateManager) AddTask(task *Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, ok := tm.tasksList[task.ID]; ok {
		log.Printf("task with ID=%s already exists", task.ID)
		return
	}

	tm.tasksList[task.ID] = task
}

func (tm *TaskStateManager) RemoveTask(task *Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, ok := tm.tasksList[task.ID]; !ok {
		log.Printf("task with ID=%s not exists", task.ID)
		return
	}

	delete(tm.tasksList, task.ID)
}

func (tm *TaskStateManager) GetTaskList() []Task {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	result := make([]Task, len(tm.tasksList))
	i := 0
	for _, val := range tm.tasksList {
		result[i] = Task{
			ID:             val.ID,
			Payload:        val.Payload,
			State:          val.State,
			CurrentRetries: val.CurrentRetries,
			MaxRetries:     val.MaxRetries,
		}
		i += 1
	}

	return result
}
