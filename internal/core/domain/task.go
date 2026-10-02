package domain

type Task struct {
	Id             ID            `json:"id"`
	Payload        string        `json:"payload"`
	MaxRetries     Attempt       `json:"max_retries"`
	CurrentRetries Attempt       `json:"current_retries"`
	State          TaskStateType `json:"state"`
}

func NewTask(id ID, payload string, maxRetries Attempt) Task {
	switch {
	case !IsIdValid(id):

	}
	return Task{
		Id:             id,
		Payload:        payload,
		MaxRetries:     maxRetries,
		CurrentRetries: 0,
		State:          TaskStateTypeCreated,
	}
}
