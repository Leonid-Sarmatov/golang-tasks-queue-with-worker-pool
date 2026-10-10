package execute

import (
	"fmt"
	"math/rand/v2"
	"time"

	"worker_pool/internal/app/conf"
	"worker_pool/internal/core/domain"
)

func Execute(task domain.Task, cfg conf.IConfig) (domain.Task, error) {

	time.Sleep(cfg.GetTaskProcessingDuration())

	if Randomizer(cfg.GetTaskProbabilityFailed()) {
		return task, domain.ErrTaskExecuteFail
	}

	fmt.Printf("task with ID=%s done\n", task.Id)

	return task, nil
}

func Randomizer(probability domain.Probability) bool {
	if !domain.IsProbabilityValid(probability) {
		return false
	}

	return rand.IntN(100) < int(probability)
}
