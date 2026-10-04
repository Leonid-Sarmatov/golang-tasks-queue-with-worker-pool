package domain

import (
	"errors"
)

var (
	ErrInvalidAttempt       = errors.New("attempt value must be non-negative")
	ErrInvalidWorkersNumber = errors.New("vorkers number must be positive")
	ErrInvalidTaskQueueSize = errors.New("task queue size must be positive")
	ErrInvalidProbability   = errors.New("probability must be in range 0-100")
	ErrInvalidId            = errors.New("ID cannot be empty")

	ErrUnknownTaskStateType = errors.New("unknown task state type")

	ErrStateMachine = errors.New("invalid transition in the state machine")

	ErrTaskNotFound      = errors.New("task with this ID not exists")
	ErrTaskAlreadyExists = errors.New("task with this ID already exists")
	ErrTaskExecuteFail   = errors.New("failed to execute task")
	ErrRetriesOwerflow   = errors.New("retries owerflow")
)
