package conf

import (
	"log"
	"os"
	"strconv"
	"time"
	"worker_pool/internal/core/domain"
)

const (
	EnvKeyWorkersNum             = "WORKERS_NUMBER"
	EnvKeyTaskQueueSize          = "TASK_QUEUE_SIZE"
	EnvKeyTaskProbabilityFailed  = "TASK_PROBABILITY_FAILED"
	EnvKeyTaskProcessingDuration = "TASK_PROCESSING_DURATION"
)

var DefaultValWorkersNum domain.WorkersNumber = 4
var DefaultValTaskQueueSize domain.TaskQueueSize = 8
var DefaultValTaskProbabilityFailed domain.Probability = 20
var DefaultValTaskProcessingDuration time.Duration = 1 * time.Second

type ConfEnv struct {
	WorkersNum             domain.WorkersNumber
	TaskQueueSize          domain.TaskQueueSize
	TaskProbabilityFailed  domain.Probability
	TaskProcessingDuration time.Duration
}

func NewConfEnv() *ConfEnv {
	var cfg ConfEnv

	cfg.WorkersNum = domain.WorkersNumber(parseInt(EnvKeyWorkersNum, int(DefaultValWorkersNum)))
	cfg.TaskQueueSize = domain.TaskQueueSize(parseInt(EnvKeyTaskQueueSize, int(DefaultValTaskQueueSize)))
	cfg.TaskProbabilityFailed = domain.Probability(parseInt(EnvKeyTaskProbabilityFailed, int(DefaultValTaskProbabilityFailed)))

	switch {
	case !domain.IsWorkersNumberValid(cfg.WorkersNum):
		cfg.WorkersNum = DefaultValWorkersNum
	case !domain.IsTaskQueueSizeValid(cfg.TaskQueueSize):
		cfg.TaskQueueSize = DefaultValTaskQueueSize
	case !domain.IsProbabilityValid(cfg.TaskProbabilityFailed):
		cfg.TaskProbabilityFailed = DefaultValTaskProbabilityFailed
	}

	return &cfg
}

func parseString(envKey string, defaultString string) string {
	val := os.Getenv(envKey)
	if val == "" {
		return defaultString
	}
	return val
}

func parseBool(envKey string, defaultBool bool) bool {
	val := os.Getenv(envKey)
	if val == "true" {
		return true
	}
	if val == "false" {
		return false
	}
	log.Printf("failed to read Bool from environments: %v, ", envKey)
	return defaultBool
}

func parseInt(envKey string, defaultInt int) int {
	intVal, err := strconv.Atoi(os.Getenv(envKey))
	if err != nil {
		log.Printf("failed to read Int32 from environments: %v, err: %v", envKey, err)
		return defaultInt
	}
	return intVal
}

func parseTimeDuration(envKey string, defaultDuration time.Duration) time.Duration {
	durationRecording, err := time.ParseDuration(os.Getenv(envKey))
	if err != nil {
		log.Printf("failed to read time.Duration from environments: %v, err: %v", envKey, err)
		return defaultDuration
	}
	return durationRecording
}
