package state

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const timeTestLayout = "2006-01-02 15:04:05"

func newTestStates() *State {
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

	t.Run("group", func(t *testing.T) {
		t.Parallel()

		states := newTestStates().All()

		for _, testCase := range testCases {

			testCase := testCase

			t.Run(fmt.Sprintf("created:%v/type:%v/final:%t", testCase.expectedCreatedAt, testCase.exptectedType, testCase.expectedIsFinal), func(t *testing.T) {
				t.Parallel()

				state := states[testCase.key]

				assert.Equal(t, testCase.expectedCreatedAt, state.CreatedAt)
				assert.Equal(t, testCase.exptectedType, state.Type)
				assert.Equal(t, testCase.expectedIsFinal, state.isFinal)
			})

		}

	})

}

func TestStatesLatest(t *testing.T) {
	t.Parallel()

	timeList := createTestTimeList()

	testCases := []struct {
		state  *State
		status *Status
	}{
		{
			state: newTestStates(),
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

	t.Run("group", func(t *testing.T) {
		t.Parallel()

		for _, testCase := range testCases {
			testCase := testCase

			t.Run(fmt.Sprintf("status:%v", testCase.status), func(t *testing.T) {
				t.Parallel()

				latest := testCase.state.Latest()

				assert.Equal(t, testCase.status, latest)
			})
		}

	})

}

func TestStateSubscribe(t *testing.T) {
	t.Parallel()

	timeList := createTestTimeList()

	testCases := []struct {
		expectedCreatedAt time.Time
		exptectedType     StatusType
		expectedIsFinal   bool
	}{
		{
			expectedCreatedAt: timeList[0],
			exptectedType:     StatusTaskStarted,
			expectedIsFinal:   false,
		}, {
			expectedCreatedAt: timeList[1],
			exptectedType:     StatusConnected,
			expectedIsFinal:   false,
		}, {
			expectedCreatedAt: timeList[2],
			exptectedType:     StatusTaskCompleted,
			expectedIsFinal:   true,
		},
	}

	t.Run("group", func(t *testing.T) {
		t.Parallel()

		states := newTestStates()
		sub, err := states.Subscribe()

		assert.NoError(t, err)

		for _, testCase := range testCases {
			testCase := testCase

			//no need run in parallel as only check state value
			t.Run(fmt.Sprintf("created:%v/type:%v/final/%t", testCase.expectedCreatedAt, testCase.exptectedType, testCase.expectedIsFinal), func(t *testing.T) {
				state := <-sub.statusCh

				assert.Equal(t, testCase.expectedCreatedAt, state.CreatedAt)
				assert.Equal(t, testCase.exptectedType, state.Type)
				assert.Equal(t, testCase.expectedIsFinal, state.isFinal)
			})
		}
	})

}

func TestStateUpdate(t *testing.T) {
	t.Parallel()

	t1, t2 := time.Now(), time.Now().Add(time.Second*10)

	testCases := []struct {
		createdAt  time.Time
		typeStatus StatusType
		isFinal    bool
	}{
		{
			createdAt:  t1,
			typeStatus: StatusTicketRegistered,
			isFinal:    false,
		}, {
			createdAt:  t2,
			typeStatus: StatusTaskRejected,
			isFinal:    true,
		},
	}

	t.Run("group", func(t *testing.T) {
		t.Parallel()

		for _, testCase := range testCases {
			testCase := testCase

			t.Run(fmt.Sprintf("created:%v/type:%v/final/%t", testCase.createdAt, testCase.typeStatus, testCase.isFinal), func(t *testing.T) {
				t.Parallel()

				states := newTestStates()

				states.Update(context.Background(), &Status{
					CreatedAt: testCase.createdAt,
					Type:      testCase.typeStatus,
					isFinal:   testCase.isFinal,
				})

				lastState := states.statuses[len(states.statuses)-1]

				assert.Equal(t, testCase.createdAt, lastState.CreatedAt)
				assert.Equal(t, testCase.typeStatus, lastState.Type)
				assert.Equal(t, testCase.isFinal, lastState.isFinal)
			})
		}
	})
}
