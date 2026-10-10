package domain

type Task struct {
	Id             ID            `json:"id"`
	Payload        string        `json:"payload"`
	MaxAttempts    Attempt       `json:"max_attempts"`
	CurrentAttempt Attempt       `json:"current_retries"`
	State          TaskStateType `json:"state"`
}

func NewTask(id ID, payload string, MaxAttempts Attempt) (Task, error) {
	switch {
	case !IsIdValid(id):
		return Task{}, ErrInvalidId
	case !IsAttemptValid(MaxAttempts):
		return Task{}, ErrInvalidAttempt
	}
	return Task{
		Id:             id,
		Payload:        payload,
		MaxAttempts:    MaxAttempts,
		CurrentAttempt: 0,
		State:          TaskStateTypeCreated,
	}, nil
}

func MakeTaskCopy(t Task) Task {
	return Task{
		Id:             t.Id,
		Payload:        t.Payload,
		MaxAttempts:    t.MaxAttempts,
		CurrentAttempt: t.CurrentAttempt,
		State:          t.State,
	}
}
