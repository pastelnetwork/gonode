package api

import (
	"context"
	"embed"
	"io"
	"net/http"
	"os"
	"time"

	artworks "github.com/pastelnetwork/walletnode/api/gen/artworks"
	artworkssvr "github.com/pastelnetwork/walletnode/api/gen/http/artworks/server"
	swaggersvr "github.com/pastelnetwork/walletnode/api/gen/http/swagger/server"
	"github.com/pastelnetwork/walletnode/api/log"
	artworksservice "github.com/pastelnetwork/walletnode/api/services/artworks"

	goahttp "goa.design/goa/v3/http"
	goahttpmiddleware "goa.design/goa/v3/http/middleware"
	goamiddleware "goa.design/goa/v3/middleware"
)

//go:embed gen/http/openapi3.json
var openapi embed.FS

//go:embed swagger
var swaggerContent embed.FS

func swaggerHandler() http.Handler {
	return http.FileServer(http.FS(swaggerContent))
}

func apiHandler() http.Handler {
	var (
		dec        = goahttp.RequestDecoder
		enc        = goahttp.ResponseEncoder
		mux        = goahttp.NewMuxer()
		errHandler = errorHandler()
	)

	artworksEndpoints := artworks.NewEndpoints(artworksservice.New())
	artworksServer := artworkssvr.New(artworksEndpoints, mux, dec, enc, errHandler, nil, artworksservice.UploadImageDecoderFunc)
	artworkssvr.Mount(mux, artworksServer)

	swaggerServer := swaggersvr.New(nil, mux, dec, enc, errHandler, nil)
	mountSwagger(mux, swaggerServer)

	servers := goahttp.Servers{
		artworksServer,
		swaggerServer,
	}
	servers.Use(goahttpmiddleware.Debug(mux, os.Stdout))

	for _, m := range artworksServer.Mounts {
		log.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range swaggerServer.Mounts {
		log.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}

	var handler http.Handler = mux

	handler = Log()(handler)
	handler = goahttpmiddleware.RequestID()(handler)

	return handler
}

func mountSwagger(mux goahttp.Muxer, server *swaggersvr.Server) {
	for _, m := range server.Mounts {
		file, err := openapi.Open(m.Method)
		if err != nil {
			continue
		}

		mux.Handle(m.Verb, m.Pattern, func(w http.ResponseWriter, r *http.Request) {
			http.ServeContent(w, r, m.Method, time.Time{}, file.(io.ReadSeeker))
		})

	}
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler() func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(goamiddleware.RequestIDKey).(string)
		_, _ = w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		log.Errorf("[%s] %s", id, err.Error())
	}
}
