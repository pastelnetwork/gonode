package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/probe"
	"github.com/pastelnetwork/gonode/probe/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// FingerprintsMethod represent Fingerprints name method
	FingerprintsMethod = "Fingerprints"

	// LoadModelsMethod represent LoadModels name method
	LoadModelsMethod = "LoadModels"
)

// Tensor implementing probe.Tensor mock for testing purpose
type Tensor struct {
	t *testing.T
	*mocks.Tensor
}

// NewMockTensor new Tensor instance
func NewMockTensor(t *testing.T) *Tensor {
	return &Tensor{
		t:      t,
		Tensor: &mocks.Tensor{},
	}
}

// ListenOnFingerprints listening Fingerprints call returns values from args
func (tensor *Tensor) ListenOnFingerprints(f probe.Fingerprints, err error) *Tensor {
	tensor.On(FingerprintsMethod, mock.Anything, mock.Anything).Return(f, err)
	return tensor
}

// ListenOnLoadModels listening LoadModels call and returns error from args
func (tensor *Tensor) ListenOnLoadModels(err error) *Tensor {
	tensor.On(LoadModelsMethod, mock.Anything).Return(err)
	return tensor
}

// AssertFingerprintsCall assertion Fingerprints call
func (tensor *Tensor) AssertFingerprintsCall(expectedCalls int, arguments ...interface{}) *Tensor {
	if expectedCalls > 0 {
		tensor.AssertCalled(tensor.t, FingerprintsMethod, arguments...)
	}
	tensor.AssertNumberOfCalls(tensor.t, FingerprintsMethod, expectedCalls)
	return tensor
}

// AssertLoadModelsCall assertion LoadModels call
func (tensor *Tensor) AssertLoadModelsCall(expectedCalls int, arguments ...interface{}) *Tensor {
	if expectedCalls > 0 {
		tensor.AssertCalled(tensor.t, LoadModelsMethod, arguments...)
	}
	tensor.AssertNumberOfCalls(tensor.t, LoadModelsMethod, expectedCalls)
	return tensor
}
