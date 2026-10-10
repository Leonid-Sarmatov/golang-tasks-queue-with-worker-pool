package kuber

import "sync/atomic"

type KubernetesAdapter struct {
	ReadyState int32
}

func NewKubernetesAdapter() *KubernetesAdapter {
	var ka KubernetesAdapter
	ka.SetStateNotReady()
	return &ka
}

func (ka *KubernetesAdapter) SetStateReadyOK() {
	atomic.SwapInt32(&ka.ReadyState, 1)
}

func (ka *KubernetesAdapter) SetStateNotReady() {
	atomic.SwapInt32(&ka.ReadyState, 0)
}

func (ka *KubernetesAdapter) IsApplicationReady() bool {
	return atomic.LoadInt32(&ka.ReadyState) == 1
}
