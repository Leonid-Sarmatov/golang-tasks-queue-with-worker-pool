package domain

type Task struct {
	Id             ID            `json:"id"`
	Payload        string        `json:"payload"`
	MaxRetries     Attempt       `json:"max_retries"`
	CurrentRetries Attempt       `json:"current_retries"`
	State          TaskStateType `json:"state"`
}

func NewTask(id ID, payload string, maxRetries Attempt) (Task, error) {
	switch {
	case !IsIdValid(id):
		return Task{}, ErrInvalidId
	case !IsAttemptValid(maxRetries):
		return Task{}, ErrInvalidAttempt
	}
	return Task{
		Id:             id,
		Payload:        payload,
		MaxRetries:     maxRetries,
		CurrentRetries: 0,
		State:          TaskStateTypeCreated,
	}, nil
}

func MakeTaskCopy(t Task) Task {
	return Task{
		Id:             t.Id,
		Payload:        t.Payload,
		MaxRetries:     t.MaxRetries,
		CurrentRetries: t.CurrentRetries,
		State:          t.State,
	}
}
