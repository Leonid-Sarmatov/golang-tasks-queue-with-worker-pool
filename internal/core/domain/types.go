package domain

type Attempt int32

type WorkersNumber int32

type TaskQueueSize int32

type Probability int32

type ID string

type TaskStateType int32

const (
	TaskStateTypeCreated TaskStateType = iota
	TaskStateTypeInQueue
	TaskStateTypeRunning
	TaskStateTypeDone
	TaskStateTypeFailed
)
