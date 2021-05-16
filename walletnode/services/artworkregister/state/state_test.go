package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const timeTestLayout = "2006-01-02 15:04:05"

func newTestState() *State {
	timeList := createTestTimeList()

	return &State{
		statuses: []*Status{
			{
				CreatedAt: timeList[0],
				Type:      StatusTaskStarted,
				isFinal:   false,
			}, {
				CreatedAt: timeList[1],
				Type:      StatusConnected,
				isFinal:   false,
			}, {
				CreatedAt: timeList[2],
				Type:      StatusTaskCompleted,
				isFinal:   true,
			},
		},
	}
}

func createTestTimeList() []time.Time {
	time1, _ := time.Parse(timeTestLayout, "2020-01-02 07:00:00")
	time2, _ := time.Parse(timeTestLayout, "2020-01-02 07:01:00")
	time3, _ := time.Parse(timeTestLayout, "2020-01-02 07:02:00")

	return []time.Time{time1, time2, time3}
}

func TestStatesAll(t *testing.T) {
	t.Parallel()

	timeList := createTestTimeList()

	testCases := []struct {
		key               int
		expectedCreatedAt time.Time
		exptectedType     StatusType
		expectedIsFinal   bool
	}{
		{
			key:               0,
			expectedCreatedAt: timeList[0],
			exptectedType:     StatusTaskStarted,
			expectedIsFinal:   false,
		}, {
			key:               1,
			expectedCreatedAt: timeList[1],
			exptectedType:     StatusConnected,
			expectedIsFinal:   false,
		}, {
			key:               2,
			expectedCreatedAt: timeList[2],
			exptectedType:     StatusTaskCompleted,
			expectedIsFinal:   true,
		},
	}

	states := newTestState().All()

	for _, testCase := range testCases {

		testCase := testCase
		state := states[testCase.key]

		assert.Equal(t, testCase.expectedCreatedAt, state.CreatedAt)
		assert.Equal(t, testCase.exptectedType, state.Type)
		assert.Equal(t, testCase.expectedIsFinal, state.isFinal)
	}

}

func TestStatesLatest(t *testing.T) {
	t.Parallel()

	timeList := createTestTimeList()

	testCases := []struct {
		state  *State
		status *Status
	}{
		{
			state: newTestState(),
			status: &Status{
				CreatedAt: timeList[2],
				Type:      StatusTaskCompleted,
				isFinal:   true,
			},
		}, {
			state:  &State{},
			status: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		latest := testCase.state.Latest()

		assert.Equal(t, testCase.status, latest)
	}

}
