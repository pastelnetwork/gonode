// Code generated by goa v3.5.3, DO NOT EDIT.
//
// artworks HTTP server
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/websocket"
	artworks "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"goa.design/plugins/v3/cors"
)

// Server lists the artworks service endpoint HTTP handlers.
type Server struct {
	Mounts            []*MountPoint
	Register          http.Handler
	RegisterTaskState http.Handler
	RegisterTask      http.Handler
	RegisterTasks     http.Handler
	UploadImage       http.Handler
	ArtSearch         http.Handler
	ArtworkGet        http.Handler
	Download          http.Handler
	CORS              http.Handler
}

// ErrorNamer is an interface implemented by generated error structs that
// exposes the name of the error as defined in the design.
type ErrorNamer interface {
	ErrorName() string
}

// MountPoint holds information about the mounted endpoints.
type MountPoint struct {
	// Method is the name of the service method served by the mounted HTTP handler.
	Method string
	// Verb is the HTTP method used to match requests to the mounted handler.
	Verb string
	// Pattern is the HTTP request path pattern used to match requests to the
	// mounted handler.
	Pattern string
}

// ArtworksUploadImageDecoderFunc is the type to decode multipart request for
// the "artworks" service "uploadImage" endpoint.
type ArtworksUploadImageDecoderFunc func(*multipart.Reader, **artworks.UploadImagePayload) error

// New instantiates HTTP handlers for all the artworks service endpoints using
// the provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *artworks.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
	upgrader goahttp.Upgrader,
	configurer *ConnConfigurer,
	artworksUploadImageDecoderFn ArtworksUploadImageDecoderFunc,
) *Server {
	if configurer == nil {
		configurer = &ConnConfigurer{}
	}
	return &Server{
		Mounts: []*MountPoint{
			{"Register", "POST", "/artworks/register"},
			{"RegisterTaskState", "GET", "/artworks/register/{taskId}/state"},
			{"RegisterTask", "GET", "/artworks/register/{taskId}"},
			{"RegisterTasks", "GET", "/artworks/register"},
			{"UploadImage", "POST", "/artworks/register/upload"},
			{"ArtSearch", "GET", "/artworks/search"},
			{"ArtworkGet", "GET", "/artworks/{txid}"},
			{"Download", "GET", "/artworks/download"},
			{"CORS", "OPTIONS", "/artworks/register"},
			{"CORS", "OPTIONS", "/artworks/register/{taskId}/state"},
			{"CORS", "OPTIONS", "/artworks/register/{taskId}"},
			{"CORS", "OPTIONS", "/artworks/register/upload"},
			{"CORS", "OPTIONS", "/artworks/search"},
			{"CORS", "OPTIONS", "/artworks/{txid}"},
			{"CORS", "OPTIONS", "/artworks/download"},
		},
		Register:          NewRegisterHandler(e.Register, mux, decoder, encoder, errhandler, formatter),
		RegisterTaskState: NewRegisterTaskStateHandler(e.RegisterTaskState, mux, decoder, encoder, errhandler, formatter, upgrader, configurer.RegisterTaskStateFn),
		RegisterTask:      NewRegisterTaskHandler(e.RegisterTask, mux, decoder, encoder, errhandler, formatter),
		RegisterTasks:     NewRegisterTasksHandler(e.RegisterTasks, mux, decoder, encoder, errhandler, formatter),
		UploadImage:       NewUploadImageHandler(e.UploadImage, mux, NewArtworksUploadImageDecoder(mux, artworksUploadImageDecoderFn), encoder, errhandler, formatter),
		ArtSearch:         NewArtSearchHandler(e.ArtSearch, mux, decoder, encoder, errhandler, formatter, upgrader, configurer.ArtSearchFn),
		ArtworkGet:        NewArtworkGetHandler(e.ArtworkGet, mux, decoder, encoder, errhandler, formatter),
		Download:          NewDownloadHandler(e.Download, mux, decoder, encoder, errhandler, formatter),
		CORS:              NewCORSHandler(),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "artworks" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.Register = m(s.Register)
	s.RegisterTaskState = m(s.RegisterTaskState)
	s.RegisterTask = m(s.RegisterTask)
	s.RegisterTasks = m(s.RegisterTasks)
	s.UploadImage = m(s.UploadImage)
	s.ArtSearch = m(s.ArtSearch)
	s.ArtworkGet = m(s.ArtworkGet)
	s.Download = m(s.Download)
	s.CORS = m(s.CORS)
}

// Mount configures the mux to serve the artworks endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountRegisterHandler(mux, h.Register)
	MountRegisterTaskStateHandler(mux, h.RegisterTaskState)
	MountRegisterTaskHandler(mux, h.RegisterTask)
	MountRegisterTasksHandler(mux, h.RegisterTasks)
	MountUploadImageHandler(mux, h.UploadImage)
	MountArtSearchHandler(mux, h.ArtSearch)
	MountArtworkGetHandler(mux, h.ArtworkGet)
	MountDownloadHandler(mux, h.Download)
	MountCORSHandler(mux, h.CORS)
}

// MountRegisterHandler configures the mux to serve the "artworks" service
// "register" endpoint.
func MountRegisterHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleArtworksOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/artworks/register", f)
}

// NewRegisterHandler creates a HTTP handler which loads the HTTP request and
// calls the "artworks" service "register" endpoint.
func NewRegisterHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeRegisterRequest(mux, decoder)
		encodeResponse = EncodeRegisterResponse(encoder)
		encodeError    = EncodeRegisterError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "register")
		ctx = context.WithValue(ctx, goa.ServiceKey, "artworks")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountRegisterTaskStateHandler configures the mux to serve the "artworks"
// service "registerTaskState" endpoint.
func MountRegisterTaskStateHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleArtworksOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/artworks/register/{taskId}/state", f)
}

// NewRegisterTaskStateHandler creates a HTTP handler which loads the HTTP
// request and calls the "artworks" service "registerTaskState" endpoint.
func NewRegisterTaskStateHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
	upgrader goahttp.Upgrader,
	configurer goahttp.ConnConfigureFunc,
) http.Handler {
	var (
		decodeRequest = DecodeRegisterTaskStateRequest(mux, decoder)
		encodeError   = EncodeRegisterTaskStateError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "registerTaskState")
		ctx = context.WithValue(ctx, goa.ServiceKey, "artworks")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		v := &artworks.RegisterTaskStateEndpointInput{
			Stream: &RegisterTaskStateServerStream{
				upgrader:   upgrader,
				configurer: configurer,
				cancel:     cancel,
				w:          w,
				r:          r,
			},
			Payload: payload.(*artworks.RegisterTaskStatePayload),
		}
		_, err = endpoint(ctx, v)
		if err != nil {
			if _, ok := err.(websocket.HandshakeError); ok {
				return
			}
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
	})
}

// MountRegisterTaskHandler configures the mux to serve the "artworks" service
// "registerTask" endpoint.
func MountRegisterTaskHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleArtworksOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/artworks/register/{taskId}", f)
}

// NewRegisterTaskHandler creates a HTTP handler which loads the HTTP request
// and calls the "artworks" service "registerTask" endpoint.
func NewRegisterTaskHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeRegisterTaskRequest(mux, decoder)
		encodeResponse = EncodeRegisterTaskResponse(encoder)
		encodeError    = EncodeRegisterTaskError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "registerTask")
		ctx = context.WithValue(ctx, goa.ServiceKey, "artworks")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountRegisterTasksHandler configures the mux to serve the "artworks" service
// "registerTasks" endpoint.
func MountRegisterTasksHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleArtworksOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/artworks/register", f)
}

// NewRegisterTasksHandler creates a HTTP handler which loads the HTTP request
// and calls the "artworks" service "registerTasks" endpoint.
func NewRegisterTasksHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		encodeResponse = EncodeRegisterTasksResponse(encoder)
		encodeError    = EncodeRegisterTasksError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "registerTasks")
		ctx = context.WithValue(ctx, goa.ServiceKey, "artworks")
		var err error
		res, err := endpoint(ctx, nil)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountUploadImageHandler configures the mux to serve the "artworks" service
// "uploadImage" endpoint.
func MountUploadImageHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleArtworksOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/artworks/register/upload", f)
}

// NewUploadImageHandler creates a HTTP handler which loads the HTTP request
// and calls the "artworks" service "uploadImage" endpoint.
func NewUploadImageHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeUploadImageRequest(mux, decoder)
		encodeResponse = EncodeUploadImageResponse(encoder)
		encodeError    = EncodeUploadImageError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "uploadImage")
		ctx = context.WithValue(ctx, goa.ServiceKey, "artworks")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountArtSearchHandler configures the mux to serve the "artworks" service
// "artSearch" endpoint.
func MountArtSearchHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleArtworksOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/artworks/search", f)
}

// NewArtSearchHandler creates a HTTP handler which loads the HTTP request and
// calls the "artworks" service "artSearch" endpoint.
func NewArtSearchHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
	upgrader goahttp.Upgrader,
	configurer goahttp.ConnConfigureFunc,
) http.Handler {
	var (
		decodeRequest = DecodeArtSearchRequest(mux, decoder)
		encodeError   = EncodeArtSearchError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "artSearch")
		ctx = context.WithValue(ctx, goa.ServiceKey, "artworks")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		v := &artworks.ArtSearchEndpointInput{
			Stream: &ArtSearchServerStream{
				upgrader:   upgrader,
				configurer: configurer,
				cancel:     cancel,
				w:          w,
				r:          r,
			},
			Payload: payload.(*artworks.ArtSearchPayload),
		}
		_, err = endpoint(ctx, v)
		if err != nil {
			if _, ok := err.(websocket.HandshakeError); ok {
				return
			}
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
	})
}

// MountArtworkGetHandler configures the mux to serve the "artworks" service
// "artworkGet" endpoint.
func MountArtworkGetHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleArtworksOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/artworks/{txid}", f)
}

// NewArtworkGetHandler creates a HTTP handler which loads the HTTP request and
// calls the "artworks" service "artworkGet" endpoint.
func NewArtworkGetHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeArtworkGetRequest(mux, decoder)
		encodeResponse = EncodeArtworkGetResponse(encoder)
		encodeError    = EncodeArtworkGetError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "artworkGet")
		ctx = context.WithValue(ctx, goa.ServiceKey, "artworks")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountDownloadHandler configures the mux to serve the "artworks" service
// "download" endpoint.
func MountDownloadHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleArtworksOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/artworks/download", f)
}

// NewDownloadHandler creates a HTTP handler which loads the HTTP request and
// calls the "artworks" service "download" endpoint.
func NewDownloadHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeDownloadRequest(mux, decoder)
		encodeResponse = EncodeDownloadResponse(encoder)
		encodeError    = EncodeDownloadError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "download")
		ctx = context.WithValue(ctx, goa.ServiceKey, "artworks")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service artworks.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleArtworksOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/artworks/register", f)
	mux.Handle("OPTIONS", "/artworks/register/{taskId}/state", f)
	mux.Handle("OPTIONS", "/artworks/register/{taskId}", f)
	mux.Handle("OPTIONS", "/artworks/register/upload", f)
	mux.Handle("OPTIONS", "/artworks/search", f)
	mux.Handle("OPTIONS", "/artworks/{txid}", f)
	mux.Handle("OPTIONS", "/artworks/download", f)
}

// NewCORSHandler creates a HTTP handler which returns a simple 200 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}

// handleArtworksOrigin applies the CORS response headers corresponding to the
// origin for the service artworks.
func handleArtworksOrigin(h http.Handler) http.Handler {
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "localhost") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			origHndlr(w, r)
			return
		}
		origHndlr(w, r)
		return
	})
}
