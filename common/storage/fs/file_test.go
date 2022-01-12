package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/stretchr/testify/assert"
)

func TestFSOpen(t *testing.T) {
	t.Parallel()

	type args struct {
		filename string
	}

	type handleFunc func(dir string, t assert.TestingT)

	testCases := []struct {
		args        args
		createfunc  handleFunc
		assertion   assert.ErrorAssertionFunc
		valueAssert assert.ValueAssertionFunc
	}{
		{
			args:        args{"test.txt"},
			assertion:   assert.NoError,
			valueAssert: assert.NotNil,
			createfunc: func(dir string, t assert.TestingT) {
				fs := &FS{
					dir: dir,
				}

				_, err := fs.Create("test.txt")
				assert.NoError(t, err)
			},
		}, {
			args:        args{"non-exit.txt"},
			assertion:   assert.Error,
			valueAssert: assert.Nil,
			createfunc:  func(dir string, t assert.TestingT) {},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run("group", func(t *testing.T) {
			dir, _ := ioutil.TempDir("", "*")
			defer os.RemoveAll(dir)

			t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
				testCase.createfunc(dir, t)
				fs := &FS{dir: dir}

				got, err := fs.Open(testCase.args.filename)
				testCase.assertion(t, err)
				testCase.valueAssert(t, got)
			})
		})

	}
}

func TestFSCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		filename string
	}

	testCases := []struct {
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			args:      args{"test-1.txt"},
			assertion: assert.NoError,
		},
	}

	t.Run("group", func(t *testing.T) {
		dir, _ := ioutil.TempDir("", "*")
		defer os.RemoveAll(dir)

		for i, testCase := range testCases {
			testCase := testCase

			t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
				fs := &FS{
					dir: dir,
				}
				got, err := fs.Create(testCase.args.filename)
				testCase.assertion(t, err)
				assert.NotNil(t, got)
				assert.FileExists(t, fmt.Sprintf("%s/%s", dir, testCase.args.filename))
			})
		}
	})
}

func TestFSRemove(t *testing.T) {
	t.Parallel()

	type args struct {
		filename string
	}

	testCases := []struct {
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			args:      args{"test-2.txt"},
			assertion: assert.NoError,
		},
	}

	t.Run("group", func(t *testing.T) {
		dir, _ := ioutil.TempDir("", "*")
		defer os.RemoveAll(dir)

		for i, testCase := range testCases {
			testCase := testCase

			t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
				fs := &FS{
					dir: dir,
				}

				_, err := fs.Create(testCase.args.filename)
				assert.NoError(t, err)

				testCase.assertion(t, fs.Remove(testCase.args.filename))
			})
		}

	})
}

func TestNewFileStorage(t *testing.T) {
	t.Parallel()

	type args struct {
		dir string
	}

	testCases := []struct {
		args args
		want storage.FileStorageInterface
	}{
		{
			args: args{"./"},
			want: &FS{dir: "./"},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			assert.Equal(t, testCase.want, NewFileStorage(testCase.args.dir))
		})
	}
}
