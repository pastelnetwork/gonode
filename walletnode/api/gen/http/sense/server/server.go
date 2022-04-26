// Code generated by goa v3.6.2, DO NOT EDIT.
//
// sense HTTP server
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design -o api/

package server

import (
	"context"
	"mime/multipart"
	"net/http"

	sense "github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"goa.design/plugins/v3/cors"
)

// Server lists the sense service endpoint HTTP handlers.
type Server struct {
	Mounts            []*MountPoint
	UploadImage       http.Handler
	StartProcessing   http.Handler
	RegisterTaskState http.Handler
	GetTaskHistory    http.Handler
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

// SenseUploadImageDecoderFunc is the type to decode multipart request for the
// "sense" service "uploadImage" endpoint.
type SenseUploadImageDecoderFunc func(*multipart.Reader, **sense.UploadImagePayload) error

// New instantiates HTTP handlers for all the sense service endpoints using the
// provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *sense.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
	upgrader goahttp.Upgrader,
	configurer *ConnConfigurer,
	senseUploadImageDecoderFn SenseUploadImageDecoderFunc,
) *Server {
	if configurer == nil {
		configurer = &ConnConfigurer{}
	}
	return &Server{
		Mounts: []*MountPoint{
			{"UploadImage", "POST", "/openapi/sense/upload"},
			{"StartProcessing", "POST", "/openapi/sense/start/{image_id}"},
			{"RegisterTaskState", "GET", "/openapi/sense/start/{taskId}/state"},
			{"GetTaskHistory", "GET", "/openapi/sense/start/{taskId}/history"},
			{"Download", "GET", "/openapi/sense/download"},
			{"CORS", "OPTIONS", "/openapi/sense/upload"},
			{"CORS", "OPTIONS", "/openapi/sense/start/{image_id}"},
			{"CORS", "OPTIONS", "/openapi/sense/start/{taskId}/state"},
			{"CORS", "OPTIONS", "/openapi/sense/start/{taskId}/history"},
			{"CORS", "OPTIONS", "/openapi/sense/download"},
		},
		UploadImage:       NewUploadImageHandler(e.UploadImage, mux, NewSenseUploadImageDecoder(mux, senseUploadImageDecoderFn), encoder, errhandler, formatter),
		StartProcessing:   NewStartProcessingHandler(e.StartProcessing, mux, decoder, encoder, errhandler, formatter),
		RegisterTaskState: NewRegisterTaskStateHandler(e.RegisterTaskState, mux, decoder, encoder, errhandler, formatter, upgrader, configurer.RegisterTaskStateFn),
		GetTaskHistory:    NewGetTaskHistoryHandler(e.GetTaskHistory, mux, decoder, encoder, errhandler, formatter),
		Download:          NewDownloadHandler(e.Download, mux, decoder, encoder, errhandler, formatter),
		CORS:              NewCORSHandler(),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "sense" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.UploadImage = m(s.UploadImage)
	s.StartProcessing = m(s.StartProcessing)
	s.RegisterTaskState = m(s.RegisterTaskState)
	s.GetTaskHistory = m(s.GetTaskHistory)
	s.Download = m(s.Download)
	s.CORS = m(s.CORS)
}

// Mount configures the mux to serve the sense endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountUploadImageHandler(mux, h.UploadImage)
	MountStartProcessingHandler(mux, h.StartProcessing)
	MountRegisterTaskStateHandler(mux, h.RegisterTaskState)
	MountGetTaskHistoryHandler(mux, h.GetTaskHistory)
	MountDownloadHandler(mux, h.Download)
	MountCORSHandler(mux, h.CORS)
}

// Mount configures the mux to serve the sense endpoints.
func (s *Server) Mount(mux goahttp.Muxer) {
	Mount(mux, s)
}

// MountUploadImageHandler configures the mux to serve the "sense" service
// "uploadImage" endpoint.
func MountUploadImageHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleSenseOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/openapi/sense/upload", f)
}

// NewUploadImageHandler creates a HTTP handler which loads the HTTP request
// and calls the "sense" service "uploadImage" endpoint.
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
		ctx = context.WithValue(ctx, goa.ServiceKey, "sense")
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

// MountStartProcessingHandler configures the mux to serve the "sense" service
// "startProcessing" endpoint.
func MountStartProcessingHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleSenseOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/openapi/sense/start/{image_id}", f)
}

// NewStartProcessingHandler creates a HTTP handler which loads the HTTP
// request and calls the "sense" service "startProcessing" endpoint.
func NewStartProcessingHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeStartProcessingRequest(mux, decoder)
		encodeResponse = EncodeStartProcessingResponse(encoder)
		encodeError    = EncodeStartProcessingError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "startProcessing")
		ctx = context.WithValue(ctx, goa.ServiceKey, "sense")
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

// MountRegisterTaskStateHandler configures the mux to serve the "sense"
// service "registerTaskState" endpoint.
func MountRegisterTaskStateHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleSenseOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/openapi/sense/start/{taskId}/state", f)
}

// NewRegisterTaskStateHandler creates a HTTP handler which loads the HTTP
// request and calls the "sense" service "registerTaskState" endpoint.
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
		ctx = context.WithValue(ctx, goa.ServiceKey, "sense")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		v := &sense.RegisterTaskStateEndpointInput{
			Stream: &RegisterTaskStateServerStream{
				upgrader:   upgrader,
				configurer: configurer,
				cancel:     cancel,
				w:          w,
				r:          r,
			},
			Payload: payload.(*sense.RegisterTaskStatePayload),
		}
		_, err = endpoint(ctx, v)
		if err != nil {
			if _, werr := w.Write(nil); werr == http.ErrHijacked {
				// Response writer has been hijacked, do not encode the error
				errhandler(ctx, w, err)
				return
			}
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
	})
}

// MountGetTaskHistoryHandler configures the mux to serve the "sense" service
// "getTaskHistory" endpoint.
func MountGetTaskHistoryHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleSenseOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/openapi/sense/start/{taskId}/history", f)
}

// NewGetTaskHistoryHandler creates a HTTP handler which loads the HTTP request
// and calls the "sense" service "getTaskHistory" endpoint.
func NewGetTaskHistoryHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeGetTaskHistoryRequest(mux, decoder)
		encodeResponse = EncodeGetTaskHistoryResponse(encoder)
		encodeError    = EncodeGetTaskHistoryError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "getTaskHistory")
		ctx = context.WithValue(ctx, goa.ServiceKey, "sense")
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

// MountDownloadHandler configures the mux to serve the "sense" service
// "download" endpoint.
func MountDownloadHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleSenseOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/openapi/sense/download", f)
}

// NewDownloadHandler creates a HTTP handler which loads the HTTP request and
// calls the "sense" service "download" endpoint.
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
		ctx = context.WithValue(ctx, goa.ServiceKey, "sense")
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
// service sense.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleSenseOrigin(h)
	mux.Handle("OPTIONS", "/openapi/sense/upload", h.ServeHTTP)
	mux.Handle("OPTIONS", "/openapi/sense/start/{image_id}", h.ServeHTTP)
	mux.Handle("OPTIONS", "/openapi/sense/start/{taskId}/state", h.ServeHTTP)
	mux.Handle("OPTIONS", "/openapi/sense/start/{taskId}/history", h.ServeHTTP)
	mux.Handle("OPTIONS", "/openapi/sense/download", h.ServeHTTP)
}

// NewCORSHandler creates a HTTP handler which returns a simple 200 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}

// HandleSenseOrigin applies the CORS response headers corresponding to the
// origin for the service sense.
func HandleSenseOrigin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOrigin(origin, "localhost") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
