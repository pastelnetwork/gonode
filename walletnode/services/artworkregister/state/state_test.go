package state

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const timeTestLayout = "2006-01-02 15:04:05"

func newTestStatus(createdAtStr string, statusType StatusType, isFinal bool) *Status {
	createdAt, _ := time.Parse(timeTestLayout, createdAtStr)

	return &Status{
		CreatedAt: createdAt,
		Type:      statusType,
		isFinal:   isFinal,
	}
}

func TestStateAll(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		state    *State
		expected []*Status
	}{
		{
			state: &State{
				statuses: []*Status{
					newTestStatus("2020-05-01 00:00:00", StatusTaskStarted, false),
					newTestStatus("2020-05-01 00:00:01", StatusTicketAccepted, false),
					newTestStatus("2020-05-01 00:01:00", StatusTaskCompleted, true),
				},
			},
			expected: []*Status{
				newTestStatus("2020-05-01 00:00:00", StatusTaskStarted, false),
				newTestStatus("2020-05-01 00:00:01", StatusTicketAccepted, false),
				newTestStatus("2020-05-01 00:01:00", StatusTaskCompleted, true),
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("test-case:%d", i), func(t *testing.T) {
			t.Parallel()

			states := testCase.state.All()
			assert.Equal(t, testCase.expected, states)
		})

	}

}

func TestStateLatest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		state  *State
		status *Status
	}{
		{
			state: &State{
				statuses: []*Status{
					newTestStatus("2020-01-02 01:00:00", StatusTaskStarted, false),
					newTestStatus("2020-01-02 01:00:02", StatusTaskCompleted, true),
				},
			},
			status: newTestStatus("2020-01-02 01:00:02", StatusTaskCompleted, true),
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("test-case:%d", i), func(t *testing.T) {
			t.Parallel()

			latest := testCase.state.Latest()
			assert.Equal(t, testCase.status, latest)
		})
	}

}

func TestStateSubscribe(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		state    *State
		update   *Status
		expected []*Status
	}{
		{
			state: &State{
				statuses: []*Status{
					newTestStatus("2020-01-02 00:00:00", StatusTaskStarted, false),
					newTestStatus("2020-01-02 00:00:05", StatusConnected, false),
				},
			},
			update: newTestStatus("2020-01-02 00:00:10", StatusTaskCompleted, true),
			expected: []*Status{
				newTestStatus("2020-01-02 00:00:00", StatusTaskStarted, false),
				newTestStatus("2020-01-02 00:00:05", StatusConnected, false),
				newTestStatus("2020-01-02 00:00:10", StatusTaskCompleted, true),
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("test-case:%d", i), func(t *testing.T) {
			t.Parallel()

			sub, err := testCase.state.Subscribe()
			assert.NoError(t, err)

			testCase.state.Update(context.Background(), testCase.update)

			var statuses []*Status
			for {
				select {
				case status := <-sub.statusCh:
					statuses = append(statuses, status)
					continue
				default:
				}
				break
			}

			assert.Equal(t, testCase.expected, statuses)
		})
	}
}

func TestStateUpdate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		states *State
		status *Status
	}{
		{
			states: &State{
				statuses: []*Status{
					newTestStatus("2021-01-01 00:00:00", StatusTaskStarted, false),
				},
			},
			status: newTestStatus("2021-01-01 00:00:10", StatusTaskRejected, true),
		},
		{
			states: &State{
				statuses: []*Status{
					newTestStatus("2021-01-01 00:00:00", StatusTaskStarted, false),
				},
			},
			status: newTestStatus("2021-01-01 00:00:10", StatusTaskCompleted, true),
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("test-case:%d", i), func(t *testing.T) {
			t.Parallel()

			states := testCase.states
			status := testCase.status
			states.Update(context.Background(), status)

			lastState := states.statuses[len(states.statuses)-1]
			assert.Equal(t, testCase.status.CreatedAt, lastState.CreatedAt)
			assert.Equal(t, testCase.status.Type, lastState.Type)
			assert.Equal(t, testCase.status.isFinal, lastState.isFinal)
		})
	}

}
