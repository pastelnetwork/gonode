package state

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const timeTestLayout = "2006-01-02 15:04:05"

func newTestStatus(createdAtStr string, status StatusType) *Status {
	createdAt, _ := time.Parse(timeTestLayout, createdAtStr)

	return &Status{
		CreatedAt: createdAt,
		Type:      status,
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
					newTestStatus("2020-01-02 00:00:00", StatusTaskStarted),
					newTestStatus("2020-01-02 00:00:03", StatusTaskCompleted),
				},
			},
			expected: []*Status{
				newTestStatus("2020-01-02 00:00:00", StatusTaskStarted),
				newTestStatus("2020-01-02 00:00:03", StatusTaskCompleted),
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			statuses := testCase.state.All()
			assert.Equal(t, testCase.expected, statuses)
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
					newTestStatus("2020-01-02 00:00:00", StatusTaskStarted),
					newTestStatus("2020-01-02 00:00:03", StatusTaskCompleted),
				},
			},
			status: newTestStatus("2020-01-02 00:00:03", StatusTaskCompleted),
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			latest := testCase.state.Latest()
			assert.Equal(t, testCase.status, latest)
		})
	}
}

func TestStateUpdate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		state    *State
		update   *Status
		expected []*Status
	}{
		{
			state: &State{
				statuses: []*Status{
					newTestStatus("2020-01-01 00:00:00", StatusTaskStarted),
				},
				updatedCh: make(chan struct{}),
			},
			update: newTestStatus("2020-01-01 00:00:40", StatusTaskCompleted),
			expected: []*Status{
				newTestStatus("2020-01-01 00:00:00", StatusTaskStarted),
				newTestStatus("2020-01-01 00:00:40", StatusTaskCompleted),
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			state := testCase.state
			state.Update(context.Background(), testCase.update)

			statuses := state.statuses
			assert.Equal(t, testCase.expected, statuses)
		})
	}
}
