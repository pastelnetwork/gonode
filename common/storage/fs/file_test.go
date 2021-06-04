package fs

import (
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/stretchr/testify/assert"
)

func TestFSOpen(t *testing.T) {
	t.Parallel()

	type fields struct {
		dir string
	}
	type args struct {
		filename string
	}

	type handleFunc func(t assert.TestingT)

	testCases := []struct {
		name        string
		fields      fields
		args        args
		createfunc  handleFunc
		removeFunc  handleFunc
		assertion   assert.ErrorAssertionFunc
		valueAssert assert.ValueAssertionFunc
	}{
		{
			name:        "valid open file",
			fields:      fields{"./"},
			args:        args{"test.txt"},
			assertion:   assert.NoError,
			valueAssert: assert.NotNil,
			createfunc: func(t assert.TestingT) {
				fs := &FS{
					dir: "./",
				}

				_, err := fs.Create("test.txt")
				assert.NoError(t, err)
			},
			removeFunc: func(t assert.TestingT) {
				fs := &FS{
					dir: "./",
				}

				assert.NoError(t, fs.Remove("test.txt"))
			},
		}, {
			name:        "invalid open file",
			fields:      fields{"./"},
			args:        args{"non-exit.txt"},
			assertion:   assert.Error,
			valueAssert: assert.Nil,
			createfunc:  func(t assert.TestingT) {},
			removeFunc:  func(t assert.TestingT) {},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run("group", func(t *testing.T) {
			testCase.createfunc(t)

			defer testCase.removeFunc(t)

			t.Run(testCase.name, func(t *testing.T) {
				fs := &FS{dir: testCase.fields.dir}

				got, err := fs.Open(testCase.args.filename)
				testCase.assertion(t, err)
				testCase.valueAssert(t, got)
			})
		})

	}
}

func TestFSCreate(t *testing.T) {
	t.Parallel()

	type fields struct {
		dir string
	}
	type args struct {
		filename string
	}

	testCases := []struct {
		name      string
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "create test-1.txt file",
			fields:    fields{"./"},
			args:      args{"test-1.txt"},
			assertion: assert.NoError,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			fs := &FS{
				dir: testCase.fields.dir,
			}
			got, err := fs.Create(testCase.args.filename)
			testCase.assertion(t, err)
			assert.NotNil(t, got)
			assert.FileExists(t, fmt.Sprintf("%s/%s", testCase.fields.dir, testCase.args.filename))
			assert.NoError(t, fs.Remove(fmt.Sprintf("%s/%s", testCase.fields.dir, testCase.args.filename)))
		})
	}
}

func TestFSRemove(t *testing.T) {
	t.Parallel()

	type fields struct {
		dir string
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FS{
				dir: tt.fields.dir,
			}
			tt.assertion(t, fs.Remove(tt.args.filename))
		})
	}
}

func TestNewFileStorage(t *testing.T) {
	t.Parallel()

	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
		want storage.FileStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewFileStorage(tt.args.dir))
		})
	}
}
