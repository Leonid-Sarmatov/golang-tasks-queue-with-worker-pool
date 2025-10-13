package core

import "sync/atomic"

type KubernetesAdapter struct {
	ReadyState int32
}

func NewKubernetesAdapter() *KubernetesAdapter {
	return &KubernetesAdapter{
		ReadyState: 0,
	}
}

func (ca *KubernetesAdapter) SetStateReadyOK() {
	atomic.SwapInt32(&ca.ReadyState, 1)
}

func (ca *KubernetesAdapter) SetStateNotReady() {
	atomic.SwapInt32(&ca.ReadyState, 0)
}

func (ca *KubernetesAdapter) IsApplicationReady() bool {
	return atomic.LoadInt32(&ca.ReadyState) == 1
}

