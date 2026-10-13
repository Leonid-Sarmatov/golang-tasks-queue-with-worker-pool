package execute

import (
	"fmt"
	"math/rand/v2"
	"time"

	"worker_pool/internal/core/domain"
)

func MakeExecutor(
	processingDuration time.Duration,
	probabilityFailed domain.Probability,
) domain.Executor {
	return func(task domain.Task) (domain.Task, error) {
		time.Sleep(processingDuration)

		if Randomizer(probabilityFailed) {
			return task, domain.ErrTaskExecuteFail
		}

		fmt.Printf("task with ID=%s done\n", task.Id)

		return task, nil
	}
}

func Randomizer(probability domain.Probability) bool {
	if !domain.IsProbabilityValid(probability) {
		return false
	}

	return rand.IntN(100) < int(probability)
}
