package conf

import (
	"time"

	"worker_pool/internal/core/domain"
)

type IConfig interface {
	GetWorkersNumber() domain.WorkersNumber
	GetTaskQueueSize() domain.TaskQueueSize
	GetTaskProbabilityFailed() domain.Probability
	GetTaskProcessingDuration() time.Duration
}
