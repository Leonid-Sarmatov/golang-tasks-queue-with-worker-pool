package domain

func ExecuteTask(task Task) (Task, error) {
	// No-op if task is already running
	if task.State == TaskStateTypeRunning {
		return task, nil
	}

	if task.State != TaskStateTypeInQueue {
		return task, ErrStateMachine
	}

	// Check retries counter
	if task.MaxRetries <= task.CurrentRetries {
		return task, ErrRetriesOwerflow
	}

	// Increment retries counter
	task.CurrentRetries += 1

	// Change state
	task.State = TaskStateTypeRunning

	return task, nil
}

func FailTask(task Task) (Task, error) {
	// No-op if task is already failed
	if task.State == TaskStateTypeFailed {
		return task, nil
	}

	if task.State != TaskStateTypeRunning {
		return task, ErrStateMachine
	}

	// Change state
	task.State = TaskStateTypeFailed

	return task, nil
}

func TaskToQueue(task Task) (Task, error) {
	// No-op if task is already in the queue
	if task.State == TaskStateTypeInQueue {
		return task, nil
	}

	if task.State != TaskStateTypeCreated && task.State != TaskStateTypeFailed {
		return task, ErrStateMachine
	}

	// Change state
	task.State = TaskStateTypeInQueue

	return task, nil
}
