package api

import (
	"embed"
	"io"
	"net/http"
	"time"

	"github.com/pastelnetwork/walletnode/api/gen/http/swagger/server"

	goahttp "goa.design/goa/v3/http"
)

//go:embed gen/http/openapi3.json
var openapi embed.FS

// Swagger represents services for swagger endpoints.
type Swagger struct{}

// Mount configures the mux to serve the swagger endpoints.
func (service *Swagger) Mount(mux goahttp.Muxer) goahttp.Server {
	srv := server.New(nil, nil, goahttp.RequestDecoder, goahttp.ResponseEncoder, errorHandler(), nil)

	for _, m := range srv.Mounts {
		file, err := openapi.Open(m.Method)
		if err != nil {
			continue
		}

		mux.Handle(m.Verb, m.Pattern, func(w http.ResponseWriter, r *http.Request) {
			http.ServeContent(w, r, m.Method, time.Time{}, file.(io.ReadSeeker))
		})

	}
	return srv
}

// NewSwagger returns the swagger Swagger implementation.
func NewSwagger() *Swagger {
	return &Swagger{}
}
