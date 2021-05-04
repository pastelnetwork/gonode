package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

const (
	shutdownTimeout = time.Second * 5
	logPrefix       = "api"
)

type Server struct {
	funcs map[string]interface{}
}

type RpcMethod struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      string   `json:"id"`
}

func New() *Server {
	return &Server{}
}

func (s *Server) Run(ctx context.Context, config *Config) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	mux := http.NewServeMux()
	//mux.Handle("/", apiHTTP)

	addr := net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port))
	srv := &http.Server{Addr: addr, Handler: mux}

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()
		log.WithContext(ctx).Infof("Server is shutting down...")

		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			errCh <- errors.Errorf("gracefully shutdown the server: %w", err)
		}
		close(errCh)
	}()

	log.WithContext(ctx).Infof("Server is ready to handle requests at %q", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Errorf("error starting server: %w", err)
	}
	defer log.WithContext(ctx).Infof("Server stoped")

	err := <-errCh
	return err

	// s.AddHandler("getinfo", s.Getinfo)

	// server := s.InitServer(":12345")

	// if err := server.ListenAndServe(); err != http.ErrServerClosed {
	// 	return fmt.Errorf("error starting server: %w", err)
	// }

	// //return server.Shutdown(ctx)

	// return nil
}

func (rpcServer *Server) AddHandler(handlerName string, handler func(method RpcMethod) ([]byte, error)) {
	if rpcServer.funcs == nil {
		rpcServer.funcs = make(map[string]interface{})
	}
	rpcServer.funcs[handlerName] = handler
}

func (rpcServer *Server) InitServer(address string) *http.Server {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var rpcMethod RpcMethod
		err := json.NewDecoder(r.Body).Decode(&rpcMethod)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if rpcServer.funcs == nil {
			err := fmt.Errorf("rpcMethod methods are not implemented  - %s", rpcMethod.Method)
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		m := reflect.ValueOf(rpcServer.funcs[strings.ToLower(rpcMethod.Method)])
		//m := reflect.ValueOf(&rpcMethod).MethodByName(strings.Title(strings.ToLower(rpcMethod.Method)))
		if m.Kind() == reflect.Invalid || m.IsNil() || m.IsZero() {
			err := fmt.Errorf("rpcMethod method not implemented  - %s", rpcMethod.Method)
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		in := make([]reflect.Value, 1)
		in[0] = reflect.ValueOf(rpcMethod)

		retValue := m.Call(in)
		if retValue == nil || len(retValue) != 2 ||
			retValue[0].Kind() != reflect.Slice ||
			retValue[1].Kind() != reflect.Interface {
			err := fmt.Errorf("invalid rpcMethod method implementation - %s", rpcMethod.Method)
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		if !retValue[1].IsNil() {
			err := retValue[1].Interface().(error)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		retBytes := retValue[0].Interface().([]byte)
		var buf bytes.Buffer

		err = json.Compact(&buf, retBytes)
		fmt.Println(err)

		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	})
	return &http.Server{Addr: address, Handler: mux}
}

func (s *Server) Getinfo(method RpcMethod) ([]byte, error) {
	type info struct {
		Method  string `json:"method"`
		PSLNode bool   `json:"psl_node"`
		Server  bool   `json:"rpc_server"`
	}

	i := info{method.Method, true, true}
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, jsonBytes); err != nil {
		fmt.Println(err)
	}
	return buffer.Bytes(), nil
}
