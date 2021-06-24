package secret

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/otrv4/ed448"
	"github.com/stretchr/testify/assert"
)

func randTestCases(start, end int) [][]byte {
	testCases := [][]byte{}
	for i := start; i <= end; i++ {
		testCase := make([]byte, int(math.Pow(2, float64(i))))
		rand.Read(testCase)
		testCases = append(testCases, testCase)
	}
	return testCases
}

func TestGobEncodeAndDecode(t *testing.T) {
	curve := ed448.NewCurve()
	_, pub, ok := curve.GenerateKeys()
	if !ok {
		t.Fatal("ed448 curve generate keys")
	}

	request := kexECDHRequest{
		PubKey: pub[:],
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(&request); err != nil {
		t.Fatalf("gob encode: %v", err)
	}

	var out kexECDHRequest
	nb := bytes.NewBuffer(buf.Bytes())
	decoder := gob.NewDecoder(nb)
	if err := decoder.Decode(&out); err != nil {
		t.Fatalf("gob decode: %v", err)
	}
	assert.Equal(t, request.PubKey, out.PubKey)
}

func handleSecretConn(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()
	defer conn.Close()

	data := make([]byte, maxPayloadSize)
	for {
		n, err := conn.Read(data)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatalf("server read: %v", err)
		}

		if _, err := conn.Write(data[:n]); err != nil {
			log.Fatalf("server write: %v", err)
		}
	}
}

func TestSecretServerAndClient(t *testing.T) {
	// start a server
	ln, err := Listen("tcp", ":8000", &Config{})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func(listener net.Listener) {
		defer wg.Done()
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			wg.Add(1)
			go handleSecretConn(&wg, conn)
		}
	}(ln)

	conn, err := Dial("tcp", "localhost:8000", &Config{})
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	for i, testCase := range randTestCases(20, 23) {
		testCase := testCase
		t.Run(fmt.Sprintf("testCase:%d", i), func(t *testing.T) {
			n, err := conn.Write(testCase)
			if err != nil {
				t.Fatalf("write: %v", err)
			}

			buf := make([]byte, n)
			m, err := conn.Read(buf)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			assert.Equal(t, testCase, buf[:m])
		})
	}

	// close the connection
	conn.Close()
	// close the listener
	ln.Close()

	wg.Wait()
}

func TestServerAndClientWithSecretWrapper(t *testing.T) {
	// start a server
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func(listener net.Listener) {
		defer wg.Done()
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			wg.Add(1)
			go handleSecretConn(&wg, Server(conn, &Config{}))
		}
	}(ln)

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	secretConn := Client(conn, &Config{})

	for i, testCase := range randTestCases(4, 5) {
		testCase := testCase
		t.Run(fmt.Sprintf("testCase:%d", i), func(t *testing.T) {
			n, err := secretConn.Write(testCase)
			if err != nil {
				t.Fatalf("write: %v", err)
			}

			buf := make([]byte, n)
			m, err := secretConn.Read(buf)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			assert.Equal(t, testCase, buf[:m])
		})
	}

	// close the connection
	conn.Close()
	// close the listener
	ln.Close()

	wg.Wait()
}

func TestSecretConnWithReadTimeout(t *testing.T) {
	// start a secret server
	ln, err := Listen("tcp", ":8000", &Config{})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func(listener net.Listener) {
		defer wg.Done()
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			wg.Add(1)
			go handleSecretConn(&wg, conn)
		}
	}(ln)

	conn, err := Dial("tcp", "localhost:8000", &Config{})
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))

	for i, testCase := range randTestCases(2, 2) {
		testCase := testCase
		t.Run(fmt.Sprintf("testCase:%d", i), func(t *testing.T) {
			n, err := conn.Write(testCase)
			if err != nil {
				t.Fatalf("write: %v", err)
			}

			buf := make([]byte, n)
			m, err := conn.Read(buf)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			assert.Equal(t, testCase, buf[:m])

			buf = make([]byte, n)
			_, err = conn.Read(buf)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "i/o timeout")
		})
	}

	// close the connection
	conn.Close()
	// close the listener
	ln.Close()

	wg.Wait()
}

func handleProtoConn(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()
	defer conn.Close()

	for {
		hdr := make([]byte, 2)
		if _, err := conn.Read(hdr); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("server read header: %v", err)
		}

		n := int(hdr[0])<<8 | int(hdr[1])

		buf := make([]byte, n)
		if _, err := conn.Read(buf); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("server read body: %v", err)
		}
	}
}

func TestSecretConnWithTCP(t *testing.T) {
	// start a secret server
	ln, err := Listen("tcp", ":8000", &Config{})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func(listener net.Listener) {
		defer wg.Done()
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			wg.Add(1)
			go handleProtoConn(&wg, conn)
		}
	}(ln)

	conn, err := Dial("tcp", "localhost:8000", &Config{})
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	for i, testCase := range randTestCases(4, 6) {
		testCase := testCase
		t.Run(fmt.Sprintf("testCase:%d", i), func(t *testing.T) {
			m := len(testCase)
			out := make([]byte, 2+m)
			out[0] = byte(m >> 8)
			out[1] = byte(m)
			copy(out[2:], testCase)
			if _, err := conn.Write(out); err != nil {
				t.Fatalf("write: %v", err)
			}
		})
	}

	// close the connection
	conn.Close()
	// close the listener
	ln.Close()

	wg.Wait()
}

func TestSecretConnWithHTTP(t *testing.T) {
	message := []byte("hello world")
	// start a secret server
	ln, err := Listen("tcp", ":8000", &Config{})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(message)
	}
	mux := http.NewServeMux()
	mux.Handle("/hi", http.HandlerFunc(handler))

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		http.Serve(ln, mux)
	}()

	tr := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return Dial("tcp", "localhost:8000", &Config{})
		},
	}
	client := http.Client{Transport: tr}
	res, err := client.Get("http://localhost:8000/hi")
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil read: %v", err)
	}
	assert.Equal(t, message, data)
}
