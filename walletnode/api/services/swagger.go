package services

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/swagger/server"

	goahttp "goa.design/goa/v3/http"
)

// Swagger represents services for swagger endpoints.
type Swagger struct {
	*Common
}

// Mount configures the mux to serve the swagger endpoints.
func (service *Swagger) Mount(_ context.Context, mux goahttp.Muxer) goahttp.Server {
	srv := server.New(nil, nil, goahttp.RequestDecoder, goahttp.ResponseEncoder, api.ErrorHandler, nil, nil)

	for _, m := range srv.Mounts {
		file, err := gen.OpenAPIContent.Open(m.Method)
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
	return &Swagger{
		Common: NewCommon(),
	}
}
