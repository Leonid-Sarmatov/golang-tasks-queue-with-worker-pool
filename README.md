# Tasks queue with worker pool and exponential backoff strategy
## Architecture
### HTTP Server and task reception
The application launches an HTTP server with endpoint for task submission:
```go
mux.HandleFunc("/enqueue", GetAddTaskHandler(wp))
```
**JSON format for new tasks**
```json
{
	"id": "123",
	"max_retries": 3
}
```
You can use this scrips for testing that: `tests/bashtests/add_one_task.sh`
### Buffered task queue
The queue is implemented using a channel with a capacity set by the `QUEUE_SIZE` environment variable
### Worker pool and graceful shutdown triggered by OS signals
```go
// Handling the termination signals
ctx, stop := signal.NotifyContext(
	context.Background(),
	os.Interrupt,    // Ctrl+C
	syscall.SIGTERM, // docker / k8s
	syscall.SIGQUIT, // debug stack dump
)
defer stop()

// ...

// Pool initialization
wp := core.NewWorkerPool(QueueSizeEnvRead(), WorkersNumEnvRead(), tsm)
wp.Run()

// ...

<-ctx.Done()
log.Printf("Received termination signal, the application will be stopped...")
wp.Stop() // Stops accepting new tasks and waits for current tasks to complete
```
### CPU intensive task emulation
```go
for range 5 {
	time.Sleep(TASK_WORK_IMITATION_DURATION * time.Millisecond)
	if Radomizer(TASK_PROBABILITY_FAILED) {
		t.State = TaskStatusFailed
		resultChan <- t.State
		log.Printf("Task ID=%v failed", t.ID)
		return
	}
	if Radomizer(50) { // 50% chance to complete early
		break
	}
}

t.State = TaskStatusDone
resultChan <- t.State
```
### Retry strategy with exponential Backoff and Jitter
```go
// Exponential backoff: 2^retries seconds
baseDelay := time.Duration(1<<uint(task.CurrentRetries)) * time.Second

// Adding Jitter: ±500 ms
jitter := time.Duration(rand.Intn(1000)-500) * time.Millisecond
totalDelay := baseDelay + jitter

log.Printf("Task ID=%s scheduled for retry in %v", task.ID, totalDelay)
time.Sleep(totalDelay)
wp.AddTastToQueue(task)
```
**Retry limit:**
```go
if task.CurrentRetries >= task.MaxRetries {
	log.Printf("Task ID=%s exceeded max retries (%d)", task.ID, task.CurrentRetries)
	// Output fails to stdout for statistics collection
	fmt.Printf("task with ID=%04d failed\n", task.ID)
	return
}
```
### Healthcheck endpoint
`GET /health` - return 200 OK if the service is alive 
The testing script for the endpoint: `tests/bashtests/health_check.sh`
## Usage and configuration
### Environment variables

| Variable                       | Default value | Description                                           |
| ------------------------------ | ------------- | ----------------------------------------------------- |
| `WORKERS`                      | 4             | Number of the workers                                 |
| `QUEUE_SIZE`                   | 64            | Capacity of the task queue channel                    |
| `TASK_PROBABILITY_FAILED`      | 25            | Probability of the task failure on each iteration (%) |
| `TASK_WORK_IMITATION_DURATION` | 100           | Delay between iterations (ms)                         |
### Run application
**Basic startup:**
```bash
go run cmd/main.go
```
**With custom configurations:**
```bash
WORKERS="10" QUEUE_SIZE="16" go run cmd/main.go > out.txt
```
**With guaranteed failures:**
```bash
WORKERS="10" TASK_PROBABILITY_FAILED="100" go run cmd/main.go > out.txt
```
## Testing and monitoring
### High load testing 
```bash
sh ./tests/bashtests/high_load_test.sh --host=localhost --port=8080 --requests-numbers=100 # Sending 100 parallel requests
```
### Task status check
```bash
go run cmd/taskstatuscheck.go
```
**Output:**
```
Task 1: ID: 26 MaxRetries: 3 CurrentRetries: 1 State: running
Task 2: ID: 4 MaxRetries: 3 CurrentRetries: 1 State: failed
Task 3: ID: 6 MaxRetries: 3 CurrentRetries: 1 State: done

...
```
### Result analysis
**Execution statistic:**
```bash
cat out.txt | sort
# Output:
task with ID=0001 done
task with ID=0002 failed
task with ID=0003 done
...

# Counting results:
cat out.txt | grep "done" | wc -l   # Successful tasks
cat out.txt | grep "failed" | wc -l # Failed tasks
```
**Validation of all tasks processed:**
```bash
echo "Total tasks processed: $(cat out.txt | wc -l)"
```
## Deploy with Minikube
### Minikube installation
You can use this script for install Minikube and necessary dependences
```bash
sh ./deploy/minikube_install.sh
```
To deploy the application, you need to build the Docker image and load it into minikube. The `run.sh` script will do this automatically
```bash
sh ./deploy/run.sh
```
### Check application pods
Monitor pod states in real-time and check their count
```sh
watch -n 1 "echo '=== PODS ==='; kubectl get pods -l app=worker-pool; echo ''; echo '=== HPA ==='; kubectl get hpa"
```
If you want to observe pods scaling, use command above and simultaneously run high-load test
```bash
sh ./tests/bashtests/high_load_test.sh --host=<node-ip> --port=30080 --requests-numbers=200
```