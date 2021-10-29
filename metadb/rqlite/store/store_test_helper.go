package store

import (
	"context"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func MustNewStoreTest(inmem bool) *Store {
	return mustNewStoreAtPaths(mustTempDir(), "", inmem, false)
}

func mustMockLister(addr string) Listener {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic("failed to create new listner")
	}
	return &mockListener{ln}
}

func mustNewStoreAtPathsLn(id, dataPath, sqlitePath string, inmem, fk bool) (*Store, net.Listener) {
	cfg := NewDBConfig(inmem)
	cfg.FKConstraints = fk
	cfg.OnDiskPath = sqlitePath

	ln := mustMockLister("localhost:0")
	s := New(context.TODO(), ln, &Config{
		DBConf: cfg,
		Dir:    dataPath,
		ID:     id,
	})
	if s == nil {
		panic("failed to create new store")
	}
	return s, ln
}

func mustNewStore(inmem bool) *Store {
	return mustNewStoreAtPaths(mustTempDir(), "", inmem, false)
}

func mustNewStoreFK(inmem bool) *Store {
	return mustNewStoreAtPaths(mustTempDir(), "", inmem, true)
}

func mustTempDir() string {
	var err error
	path, err := ioutil.TempDir("", "rqlilte-test-")
	if err != nil {
		panic("failed to create temp dir")
	}
	return path
}

func mustNewStoreSQLitePath() (*Store, string) {
	dataDir := mustTempDir()
	sqliteDir := mustTempDir()
	sqlitePath := filepath.Join(sqliteDir, "explicit-path.db")
	return mustNewStoreAtPaths(dataDir, sqlitePath, false, true), sqlitePath
}

func mustNewStoreAtPaths(dataPath, sqlitePath string, inmem, fk bool) *Store {
	s, _ := mustNewStoreAtPathsLn(randomString(), dataPath, sqlitePath, inmem, fk)
	return s
}

func randomString() string {
	var output strings.Builder
	chars := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
	for i := 0; i < 20; i++ {
		random := rand.Intn(len(chars))
		randomChar := chars[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

type mockListener struct {
	ln net.Listener
}

func (m *mockListener) Dial(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("tcp", addr, timeout)
}

func (m *mockListener) Accept() (net.Conn, error) { return m.ln.Accept() }

func (m *mockListener) Close() error { return m.ln.Close() }

func (m *mockListener) Addr() net.Addr { return m.ln.Addr() }

func mustWriteFile(path, contents string) {
	err := os.WriteFile(path, []byte(contents), 0644)
	if err != nil {
		panic("failed to write to file")
	}
}

func mustParseDuration(t string) time.Duration {
	d, err := time.ParseDuration(t)
	if err != nil {
		panic("failed to parse duration")
	}
	return d
}
