# Task queue with a worker pool and exponential backoff strategy

## Architecture

### HTTP Server and task reception
The application launches an HTTP server with an endpoint for task submission:
```go
mux.HandleFunc("/enqueue", GetAddTaskHandler(wp))
```
**JSON format for new tasks**
```json
{
	"id":          "123",
    "payload":     "any text",
	"max_retries": 3
}
```
Use this scrips to test task submission:
```bash
sh ./tests/bashtests/add_one_task.sh 
```

### Buffered task queue
The queue is implemented using a channel with a capacity set by the `TASK_QUEUE_SIZE` environment variable.

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
wp := workerpool.NewWorkerPool(
    ctx,
    cfg.GetWorkersNumber(),
    cfg.GetTaskQueueSize(),
    execute_func)
wp.Run()

// ...

<-ctx.Done()
log.Printf("Received termination signal, the application will be stopped...")

wp.Stop() // Stops accepting new tasks and waits for current tasks to complete
```

### Task processing simulation
```go
time.Sleep(processingDuration)

if Randomizer(probabilityFailed) {
    return task, domain.ErrTaskExecuteFail
}

fmt.Printf("task with ID=%s done\n", task.Id)

return task, nil
```

### Retry strategy with exponential Backoff and Jitter
```go
// Calculate exponential backoff delay (2^retries seconds)
x := time.Duration(1<<uint(t.CurrentAttempt)) * time.Second

// Add random jitter to delay to spread retries (±500 ms)
x += time.Duration(rand.Intn(1000)-500) * time.Millisecond
```
**Retry limit:**
```go
func ExecuteTask(task Task) (Task, error) {
	// No-op if task is already running
	if task.State == TaskStateTypeRunning {
		return task, nil
	}

	if task.State != TaskStateTypeInQueue {
		return task, ErrStateMachine
	}

	// Check attempt counter
	if task.MaxAttempts <= task.CurrentAttempt {
		return task, ErrRetriesOwerflow
	}

	// Increment attempt counter
	task.CurrentAttempt += 1

	// Change state
	task.State = TaskStateTypeRunning

	return task, nil
}
```

## Usage and configuration

### Environment variables

| Variable                       | Default value | Description                                           |
| ------------------------------ | ------------- | ----------------------------------------------------- |
| `WORKERS_NUMBER`               | 4             | Number of the workers                                 |
| `TASK_QUEUE_SIZE`              | 8             | Capacity of the task queue channel                    |
| `TASK_PROBABILITY_FAILED`      | 20            | Probability of the task failure per execution attempt (%) |
| `TASK_PROCESSING_DURATION`     | 1s            | Delay in the execution function                       |
| `HTTP_PORT`                    | :8080         | Port on which the HTTP server will be launched                      |
| `HTTP_READ_TIMEOUT`            | 1s            | HTTP read timeout                                     |
| `HTTP_WRITE_TIMEOUT`           | 1s            | HTTP write timeout                                    |
| `HTTP_SHUTDOWN_TIMEOUT`        | 1s            | Durarion for graceful shutdown                        |

### Run application
**Basic launching:**
```bash
go run cmd/main.go
```
**With custom configurations:**
```bash
WORKERS_NUMBER="10" TASK_QUEUE_SIZE="16" go run cmd/main.go > out.txt
```
**With guaranteed failures:**
```bash
WORKERS_NUMBER="10" TASK_PROCESSING_DURATION="100" go run cmd/main.go > out.txt
```

## Deploy with Minikube

### Minikube installation
Use this sript to install Minikube and required dependences.
```bash
sh ./deploy/minikube_install.sh
```
To deploy the application, you need to build the Docker image and load it into Minikube. The `run.sh` script will do this automatically.
```bash
sh ./deploy/run.sh
```

### Observe pod autoscaling
Monitor pod states and check their count in real time:
```sh
watch -n 1 "echo '=== PODS ==='; kubectl get pods -l app=worker-pool; echo ''; echo '=== HPA ==='; kubectl get hpa"
```
If you want to observe pods scaling, run the command above and run the high-load test in another terminal:
```bash
sh ./tests/bashtests/high_load_test.sh --port=30080 --requests-numbers=200
```

### Check graceful shutdown logs
Open log messages stream from all pods with command:
```bash
kubectl logs -n default -l app=worker-pool --all-containers=true --prefix=true --timestamps=true --tail=50 --max-log-requests=10 -f
```
Next you should delete any pod using the command:
```bash
kubectl delete pod <pod-name> --wait=false 
```
You can find pod names via:
```bash
kubectl get pods
```

## Testing and monitoring

### Health check
`GET /health` - return 200 OK if the service is alive. The the endpoint with:
```bash
sh ./tests/bashtests/health_check.sh 
```

### High-load testing 
```bash
sh ./tests/bashtests/high_load_test.sh --port=8080 --requests-numbers=100 # Sending 100 parallel requests
```