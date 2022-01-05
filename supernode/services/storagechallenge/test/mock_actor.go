package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/messaging/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	SendMethod = "Send"
)

type Actor struct {
	t *testing.T
	*mocks.Actor
}

func (m *Actor) ListenOnSend(err error) *Actor {
	m.On(SendMethod, mock.Anything, mock.Anything, mock.Anything).Return(err)
	return m
}

func (m *Actor) AssertSendCall(expectedCalls int, arguments ...interface{}) *Actor {
	if expectedCalls > 0 {
		m.AssertCalled(m.t, SendMethod, arguments...)
	}

	m.AssertNumberOfCalls(m.t, SendMethod, expectedCalls)
	return m
}

func NewMockActor(t *testing.T) *Actor {
	return &Actor{
		t:     t,
		Actor: &mocks.Actor{},
	}
}
