package core

import (
	"fmt"
	//"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	TaskStatusCreated = "created"
	TaskStatusQueued  = "queued"
	TaskStatusRunning = "running"
	TaskStatusDone    = "done"
	TaskStatusFailed  = "failed"
)

var (
	TASK_PROBABILITY_FAILED      = 25
	TASK_WORK_IMITATION_DURATION = 100
)

type Task struct {
	ID             string `json:"id"`
	Payload        string `json:"payload"`
	MaxRetries     int    `json:"max_retries"`
	CurrentRetries int    `json:"current_retries"`
	State          string `json:"state"`
}

func NewTask(id string, pld string, mr int) *Task {
	return &Task{
		ID:             id,
		Payload:        pld,
		MaxRetries:     mr,
		CurrentRetries: 0,
		State:          TaskStatusCreated,
	}
}

func (t *Task) Execute(resultChan chan string) {
	// Меняем статус задачи
	t.State = TaskStatusRunning
	// Инкрементируем счетчик попыток решить задачу
	t.CurrentRetries += 1
	//log.Printf("<task.go> task with ID=%v running, reties=%d (max %d)", t.ID, t.CurrentRetries, t.MaxRetries)

	// Задача может выполняться 100, 200, 300, 400 или 500 мс, в зависимости от рандомайзера
	for range 5 {
		// Спим 100 мс
		time.Sleep(time.Duration(int64(TASK_WORK_IMITATION_DURATION) * int64(time.Millisecond)))
		// Решаем: задача упала или нет
		if Radomizer(TASK_PROBABILITY_FAILED) {
			// Задача упала, меняем статус
			t.State = TaskStatusFailed
			resultChan <- t.State
			//log.Printf("<task.go> task with ID=%v failed", t.ID)
			return
		}
		// Реашем: задача завершилась, или еще потянем время
		if Radomizer(50) {
			break
		}
	}

	// Завершаем задачу
	t.State = TaskStatusDone
	resultChan <- t.State
	//log.Printf("<task.go> task with ID=%v done", t.ID)

	// Запись в стандартный вывод для аккумуляции завершенных задач в файл (go run main.go > out.txt)
	if val, err := strconv.Atoi(t.ID); err == nil {
		fmt.Printf("task with ID=%04d done\n", val)
	}
}

func Radomizer(probability int) bool {
	if probability <= 0 {
		return false
	}
	if probability >= 100 {
		return true
	}

	rand.Seed(time.Now().UnixNano())
	randomValue := rand.Intn(100) + 1

	return randomValue <= probability
}
