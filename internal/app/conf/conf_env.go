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

	EnvKeyHttpPort            = "HTTP_PORT"
	EnvKeyHttpReadTimeout     = "HTTP_READ_TIMEOUT"
	EnvKeyHttpWriteTimeout    = "HTTP_WRITE_TIMEOUT"
	EnvKeyHttpShutdownTimeout = "HTTP_SHUTDOWN_TIMEOUT"
)

var DefaultValWorkersNum domain.WorkersNumber = 4
var DefaultValTaskQueueSize domain.TaskQueueSize = 8
var DefaultValTaskProbabilityFailed domain.Probability = 20
var DefaultValTaskProcessingDuration time.Duration = 1 * time.Second

var DefaultHttpPort string = ":8080"
var DefaultHttpReadTimeout time.Duration = 1 * time.Second
var DefaultHttpWriteTimeout time.Duration = 1 * time.Second
var DefaultHttpShutdownTimeout time.Duration = 5 * time.Second

type ConfEnv struct {
	WorkersNum             domain.WorkersNumber
	TaskQueueSize          domain.TaskQueueSize
	TaskProbabilityFailed  domain.Probability
	TaskProcessingDuration time.Duration

	HttpPort            string
	HttpReadTimeout     time.Duration
	HttpWriteTimeout    time.Duration
	HttpShutdownTimeout time.Duration
}

func NewConfEnv() IConfig {
	var cfg ConfEnv

	cfg.WorkersNum = domain.WorkersNumber(parseInt(EnvKeyWorkersNum, int(DefaultValWorkersNum)))
	cfg.TaskQueueSize = domain.TaskQueueSize(parseInt(EnvKeyTaskQueueSize, int(DefaultValTaskQueueSize)))
	cfg.TaskProbabilityFailed = domain.Probability(parseInt(EnvKeyTaskProbabilityFailed, int(DefaultValTaskProbabilityFailed)))
	cfg.TaskProcessingDuration = parseTimeDuration(EnvKeyTaskProcessingDuration, DefaultValTaskProcessingDuration)

	cfg.HttpPort = parseString(EnvKeyHttpPort, DefaultHttpPort)
	cfg.HttpReadTimeout = parseTimeDuration(EnvKeyHttpReadTimeout, DefaultHttpReadTimeout)
	cfg.HttpWriteTimeout = parseTimeDuration(EnvKeyHttpWriteTimeout, DefaultHttpWriteTimeout)
	cfg.HttpShutdownTimeout = parseTimeDuration(EnvKeyHttpShutdownTimeout, DefaultHttpShutdownTimeout)

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

func (cfg *ConfEnv) GetWorkersNumber() domain.WorkersNumber {
	return cfg.WorkersNum
}

func (cfg *ConfEnv) GetTaskQueueSize() domain.TaskQueueSize {
	return cfg.TaskQueueSize
}

func (cfg *ConfEnv) GetTaskProbabilityFailed() domain.Probability {
	return cfg.TaskProbabilityFailed
}

func (cfg *ConfEnv) GetTaskProcessingDuration() time.Duration {
	return cfg.TaskProcessingDuration
}

func (cfg *ConfEnv) GetHttpPort() string {
	return cfg.HttpPort
}

func (cfg *ConfEnv) GetHttpReadTimeout() time.Duration {
	return cfg.HttpReadTimeout
}

func (cfg *ConfEnv) GetHttpWriteTimeout() time.Duration {
	return cfg.HttpWriteTimeout
}

func (cfg *ConfEnv) GetHttpShutdownTimeout() time.Duration {
	return cfg.HttpShutdownTimeout
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
