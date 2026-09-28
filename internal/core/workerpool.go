package core

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var (
	ErrQueueIsFull   = errors.New("queue is full")
	ErrQueueIsLocked = errors.New("worker pool has stopped, queue is locked")
	ErrInvalidData   = errors.New("invalid data")
)

type WorkerPool struct {
	TaskQueue        chan *Task
	WorkersNum       int
	WorkerLocker     chan struct{}
	Shutdown         chan struct{}
	wg               sync.WaitGroup
	taskStateManager *TaskStateManager
}

func NewWorkerPool(qsize, wn int, tsm *TaskStateManager) *WorkerPool {
	var wp WorkerPool

	// Входная очередь задач
	wp.TaskQueue = make(chan *Task, qsize)

	// Канал-сигнал для корректного завершения работы
	wp.Shutdown = make(chan struct{}, wn)

	// канал контроля заполнения воркеров
	wp.WorkerLocker = make(chan struct{}, wn)

	// Создание группы для ожидания завершения всех горутин
	wp.wg = sync.WaitGroup{}

	wp.taskStateManager = tsm

	return &wp
}

func (wp *WorkerPool) AddTastToQueue(task *Task) error {
	// Количество повторов и максимальное количество повторов не могут быть отрицательными
	if task.CurrentRetries < 0 || task.MaxRetries < 0 {
		return ErrInvalidData
	}
	// Если программа завершает работу, канал входной очереди будет закрыт
	// для избежания паники из за записи в закрытый канал, очередь надо блокировать
	select {
	case <-wp.Shutdown:
		return ErrQueueIsLocked
	default:
		// Очередь доступна для записи
		wp.TaskQueue <- task
	}
	return nil
}

func (wp *WorkerPool) Run() {
	// Горутина, которая прослушивает очередь входящих сообщений
	go func() {
		// При закрытии канала, обрабатываем оставшиеся задачи
		for inputTask := range wp.TaskQueue {
			// Если все воркеры заняты, мы заблокируемся до освобождения хотя бы одного
			wp.WorkerLocker <- struct{}{}
			//log.Printf("<worckerpool.go> received the task from the queue, ID=%v", inputTask.ID)
			// запуск логики процесса обработки задачи
			wp.taskProcess(inputTask)
		}
	}()
}

func (wp *WorkerPool) taskProcess(task *Task) {
	// Создание канала для результата выполнения задачи
	resultChan := make(chan string, 1)

	// Горутина выполнение задачи
	wp.wg.Add(1)
	wp.taskStateManager.AddTask(task)
	go func() {
		// Уменьшаем счетчик запущеных задач
		defer wp.wg.Done()
		// По завершению освобождаем воркер
		defer func() {
			<-wp.WorkerLocker
		}()
		// По завершению удаляем задачу из списка обрабатываемых задач
		defer wp.taskStateManager.RemoveTask(task)
		// Выполнение задачи
		task.Execute(resultChan)
	}()

	// Горутина чтения канала результата выполнения задачи
	go func() {
		select {
		// Ждем результат
		case res := <-resultChan:
			// Закрываем более не нужный канал
			close(resultChan)
			// Если задача упала, отправляем ее снова в очередь с задержкой
			if res == TaskStatusFailed {
				// Проверка количества повторов
				if task.CurrentRetries >= task.MaxRetries {
					// Повторов слишком много, на повтор не отправляем
					//log.Printf("<worckerpool.go> task with ID=%s reties overflow (reties = %d), task will be ignore", task.ID, task.CurrentRetries)
					// Запись в стандартный вывод для аккумуляции проваленых задач в файл (go run main.go > fail.txt)
					if val, err := strconv.Atoi(task.ID); err == nil {
						fmt.Printf("task with ID=%04d failed\n", val)
					}
					return
				}
				// Реализуем экспоненциальный бэкофф (2^повторение + 1 секунда)
				x := time.Duration(1<<uint(task.CurrentRetries)) * time.Second
				//x := time.Duration(1 * time.Second)
				// Добавляем джиттер (разброс от -500 до +500 мс ко времени задержки)
				x += time.Duration(rand.Intn(1000)-500) * time.Millisecond
				// Ждем заданное время и добавляем задачу обратно в очередь
				//log.Printf("<worckerpool.go> task with ID=%s will be returned in queue after %d ms", task.ID, x)
				time.Sleep(time.Duration(x))
				wp.AddTastToQueue(task)
			}
		case <-wp.Shutdown:
			return
		}
	}()
}

func (wp *WorkerPool) Stop() {
	// Закрываем канал блокировки очереди задач
	close(wp.Shutdown)
	// Закрываем канал очереди задач
	close(wp.TaskQueue)
	// Ждем завершения выполнения запущеных задач (отложенные задачи будут отброшены!!!)
	log.Printf("<worckerpool.go> waiting for tasks that haven't completed yet, reties tasks will be ignore!")
	wp.wg.Wait()
	log.Printf("<worckerpool.go> all running tasks was done, the worker pool will be shutdown")
	// Закрываем канал счетчика(блокировщика) количества доступных воркеров
	close(wp.WorkerLocker)
}
