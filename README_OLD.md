# Структура программы
## Прием задач
Для приема заданий приложение поднимает HTTP сервер с эндпоинтом `/enqueue`
```go
mux.HandleFunc("/enqueue", GetAddTaskHandler(wp))
```
## JSON с новой задачей
Эндпоинт приема заданий декодирует JSON средствами стандартной библиотеки. Для теста можно воспользоваться bash-скриптом `add_one_task.sh` из директории `/tests/bashtests`
## Буферизированная очередь
Очередь задачь реализована при помощи канала, емкость которого задается при помощи переменной окружения
## Обработка пуллом воркеров и завершение программы по сигналам от OS
Пулл воркеров предусматривает завершение работы с ожиданием выполнения незавершенных задач:
```go
// Создаем пулл воркеров
wp := core.NewWorkerPool(QueueSizeEnvRead(), WorkersNumEnvRead(), tsm)
wp.Run()

// ...

// Создаем канал с сигналом об остановки сервиса
osSignalsChan := make(chan os.Signal, 1)
signal.Notify(osSignalsChan, os.Interrupt, syscall.SIGTERM)

// Ждем сигнал об остановке
<-osSignalsChan
log.Printf("A SIGINT or SIGTERM signal is received, and the application will be stopped...")
wp.Stop()
```
Программа перестает принимать новые задачи в очередь, при этом доделывает то, что уже запущено и то, что осталось в очереди. 
## Обработка задачи
Симуляция работы тяжелой задачи состоит из нескольких итераций цикла. В начале каждой итерации стоит `time.Sleep`, после которого задача может завершиться успешно (с заданной вероятностью), либо "упасть" (тоже с заданной вероятностью). В ином случае происходит переход на следующую итерацию, по завершению всех итераций задача завершается успешно
```go
// Задача может выполняться 100, 200, 300, 400 или 500 мс, в зависимости от рандомайзера
for range 5 {
	// Спим 100 мс
	time.Sleep(time.Duration(int64(TASK_WORK_IMITATION_DURATION) * int64(time.Millisecond)))
	// Решаем: задача упала или нет
	if Radomizer(TASK_PROBABILITY_FAILED) {
		// Задача упала, меняем статус
		t.State = TaskStatusFailed
		resultChan <- t.State
		log.Printf("<task.go> task with ID=%v failed", t.ID)
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
log.Printf("<task.go> task with ID=%v done", t.ID)
```
## Экспонециальный бэкофф с джиттером
Задачи которые "провалились" отправляются в конец очередь, что бы не занимать воркер. При таком подходе воркер сразу освобождается и принимается за следующую задачу из очереди. Задача хранит в себе счетчик повторов, на основе которого можно вычислить задержку и выполнить `time.Sleep` перед тем, как обратно положить задачу в очередь.
```go
// Реализуем экспоненциальный бэкофф (2^счетчик + 1)
x := time.Duration(1<<uint(task.CurrentRetries)) * time.Second
// Добавляем джиттер (разброс от -500 до +500 мс ко времени задержки)
x += time.Duration(rand.Intn(1000)-500) * time.Millisecond
// Ждем заданное время и добавляем задачу обратно в очередь
log.Printf("<worckerpool.go> task with ID=%s will be returned in queue after %d ms", task.ID, x)
time.Sleep(time.Duration(x))
wp.AddTastToQueue(task)
```
Если счетчик превышает макимальное значение для этой задачи, то задача отбрасывается как "нерешаемая", о чем сообщается в стандартный вывод:
```go
// Проверка количества повторов
if task.CurrentRetries >= task.MaxRetries {
	// Повторов слишком много, на повтор не отправляем
	log.Printf("<worckerpool.go> task with ID=%s reties overflow (reties = %d), task will be ignore", task.ID, task.CurrentRetries)
	// Запись в стандартный вывод для аккумуляции проваленых задач в файл (go run main.go > fail.txt)
	if val, err := strconv.Atoi(task.ID); err == nil {
		fmt.Printf("task with ID=%04d failed\n", val)
	}
	return
}
```
## Эндпоинт healthcheck
Написан простейший HTTP хэндлер, выдающий статус 200 и строку о том что все впорядке. Протестировать можно при помощи скрипта `health_check.sh`
# Запуск и тестирование
****
**Запуск без переменных окружения**
Точка входа лежит в директории `/cmd`. Если не задать переменные окружения, то будут выбраны значения по умолчанию
```bash
go run main.go
```
Вывод:
```bash
2025/09/07 22:53:55 queue size not set! selected default value 64
2025/09/07 22:53:55 workers not set! selected default value 4
2025/09/07 22:53:55 <server.go> http server successful started!
```
****
**Запуск с переменными окружения**
```bash
WORKERS="10" QUEUE_SIZE="16" go run main.go > out.txt
```
Для упрощения тестирования предусмотрен вывод некоторых данных не только в лог, а еще и в стандартный вывод Linux. В файл `out.txt` будет отправлена информация о задачах которые в итоге получили статус `done` и не были завершены и остались в статусе `failed`
****
Тест с заданной вероятностью "падения" задачи . Если ее выкрутить на 100, то все задачи завершатся с `failed`
```go
WORKERS="10" TASK_PROBABILITY_FAILED="100" go run main.go > out.txt
```
## Нагрузочное тестирование
Отправка 100 задач одновременно:
```bash
sh ./high_load_test.sh
```
## Проверка состояний задач
```bash
go run taskstatuschech.go
```
Пример вывода:
```bash
Task   1:       ID: 26          MaxRetries:   3         CurrentRetries:   1     State: running
Task   2:       ID: 4           MaxRetries:   3         CurrentRetries:   1     State: failed
Task   3:       ID: 6           MaxRetries:   3         CurrentRetries:   1     State: done
Task   4:       ID: 5           MaxRetries:   3         CurrentRetries:   1     State: done
Task   5:       ID: 8           MaxRetries:   3         CurrentRetries:   1     State: done
Task   6:       ID: 13          MaxRetries:   3         CurrentRetries:   1     State: done
Task   7:       ID: 35          MaxRetries:   3         CurrentRetries:   1     State: done
Task   8:       ID: 36          MaxRetries:   3         CurrentRetries:   1     State: running
Task   9:       ID: 1           MaxRetries:   3         CurrentRetries:   1     State: done
...
```
## Анализ вывода
```bash
cat out.txt | sort
```
Примерный вывод:
```bash
task with ID=0001 done
task with ID=0002 failed
task with ID=0003 failed
task with ID=0004 done
task with ID=0005 failed
task with ID=0006 done
task with ID=0007 done
task with ID=0008 done
...
```
Для того что бы убедиться что все задачи были обработаны, можно подсчитать количество строк:
```bash
...
task with ID=0096 failed
task with ID=0097 done
task with ID=0098 done
task with ID=0099 done
task with ID=0100 done
[anyuser@archlinux cmd]$ cat out.txt | sort | wc -l
100
[anyuser@archlinux cmd]$ 
```