// Code generated by goa v3.5.3, DO NOT EDIT.
//
// userdatas HTTP server
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"context"
	"mime/multipart"
	"net/http"

	userdatas "github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"goa.design/plugins/v3/cors"
)

// Server lists the userdatas service endpoint HTTP handlers.
type Server struct {
	Mounts         []*MountPoint
	CreateUserdata http.Handler
	UpdateUserdata http.Handler
	GetUserdata    http.Handler
	CORS           http.Handler
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

// UserdatasCreateUserdataDecoderFunc is the type to decode multipart request
// for the "userdatas" service "createUserdata" endpoint.
type UserdatasCreateUserdataDecoderFunc func(*multipart.Reader, **userdatas.CreateUserdataPayload) error

// UserdatasUpdateUserdataDecoderFunc is the type to decode multipart request
// for the "userdatas" service "updateUserdata" endpoint.
type UserdatasUpdateUserdataDecoderFunc func(*multipart.Reader, **userdatas.UpdateUserdataPayload) error

// New instantiates HTTP handlers for all the userdatas service endpoints using
// the provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *userdatas.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
	userdatasCreateUserdataDecoderFn UserdatasCreateUserdataDecoderFunc,
	userdatasUpdateUserdataDecoderFn UserdatasUpdateUserdataDecoderFunc,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"CreateUserdata", "POST", "/userdatas/create"},
			{"UpdateUserdata", "POST", "/userdatas/update"},
			{"GetUserdata", "GET", "/userdatas/{pastelid}"},
			{"CORS", "OPTIONS", "/userdatas/create"},
			{"CORS", "OPTIONS", "/userdatas/update"},
			{"CORS", "OPTIONS", "/userdatas/{pastelid}"},
		},
		CreateUserdata: NewCreateUserdataHandler(e.CreateUserdata, mux, NewUserdatasCreateUserdataDecoder(mux, userdatasCreateUserdataDecoderFn), encoder, errhandler, formatter),
		UpdateUserdata: NewUpdateUserdataHandler(e.UpdateUserdata, mux, NewUserdatasUpdateUserdataDecoder(mux, userdatasUpdateUserdataDecoderFn), encoder, errhandler, formatter),
		GetUserdata:    NewGetUserdataHandler(e.GetUserdata, mux, decoder, encoder, errhandler, formatter),
		CORS:           NewCORSHandler(),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "userdatas" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.CreateUserdata = m(s.CreateUserdata)
	s.UpdateUserdata = m(s.UpdateUserdata)
	s.GetUserdata = m(s.GetUserdata)
	s.CORS = m(s.CORS)
}

// Mount configures the mux to serve the userdatas endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountCreateUserdataHandler(mux, h.CreateUserdata)
	MountUpdateUserdataHandler(mux, h.UpdateUserdata)
	MountGetUserdataHandler(mux, h.GetUserdata)
	MountCORSHandler(mux, h.CORS)
}

// MountCreateUserdataHandler configures the mux to serve the "userdatas"
// service "createUserdata" endpoint.
func MountCreateUserdataHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleUserdatasOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/userdatas/create", f)
}

// NewCreateUserdataHandler creates a HTTP handler which loads the HTTP request
// and calls the "userdatas" service "createUserdata" endpoint.
func NewCreateUserdataHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeCreateUserdataRequest(mux, decoder)
		encodeResponse = EncodeCreateUserdataResponse(encoder)
		encodeError    = EncodeCreateUserdataError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "createUserdata")
		ctx = context.WithValue(ctx, goa.ServiceKey, "userdatas")
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

// MountUpdateUserdataHandler configures the mux to serve the "userdatas"
// service "updateUserdata" endpoint.
func MountUpdateUserdataHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleUserdatasOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/userdatas/update", f)
}

// NewUpdateUserdataHandler creates a HTTP handler which loads the HTTP request
// and calls the "userdatas" service "updateUserdata" endpoint.
func NewUpdateUserdataHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeUpdateUserdataRequest(mux, decoder)
		encodeResponse = EncodeUpdateUserdataResponse(encoder)
		encodeError    = EncodeUpdateUserdataError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "updateUserdata")
		ctx = context.WithValue(ctx, goa.ServiceKey, "userdatas")
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

// MountGetUserdataHandler configures the mux to serve the "userdatas" service
// "getUserdata" endpoint.
func MountGetUserdataHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleUserdatasOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/userdatas/{pastelid}", f)
}

// NewGetUserdataHandler creates a HTTP handler which loads the HTTP request
// and calls the "userdatas" service "getUserdata" endpoint.
func NewGetUserdataHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeGetUserdataRequest(mux, decoder)
		encodeResponse = EncodeGetUserdataResponse(encoder)
		encodeError    = EncodeGetUserdataError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "getUserdata")
		ctx = context.WithValue(ctx, goa.ServiceKey, "userdatas")
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
// service userdatas.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleUserdatasOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/userdatas/create", f)
	mux.Handle("OPTIONS", "/userdatas/update", f)
	mux.Handle("OPTIONS", "/userdatas/{pastelid}", f)
}

// NewCORSHandler creates a HTTP handler which returns a simple 200 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}

// handleUserdatasOrigin applies the CORS response headers corresponding to the
// origin for the service userdatas.
func handleUserdatasOrigin(h http.Handler) http.Handler {
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
