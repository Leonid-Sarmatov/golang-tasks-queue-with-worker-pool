package domain

func IsAttemptValid(a Attempt) bool {
	return 0 <= a
}

func IsWorkersNumberValid(wn WorkersNumber) bool {
	return 0 < wn
}

func IsTaskQueueSizeValid(tqs TaskQueueSize) bool {
	return 0 < tqs
}

func IsProbabilityValid(p Probability) bool {
	return 0 <= p && p <= 100
}

func IsIdValid(id ID) bool {
	return id != ""
}
