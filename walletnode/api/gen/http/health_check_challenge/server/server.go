// Code generated by goa v3.15.0, DO NOT EDIT.
//
// HealthCheckChallenge HTTP server
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"context"
	"net/http"

	healthcheckchallenge "github.com/pastelnetwork/gonode/walletnode/api/gen/health_check_challenge"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"goa.design/plugins/v3/cors"
)

// Server lists the HealthCheckChallenge service endpoint HTTP handlers.
type Server struct {
	Mounts          []*MountPoint
	GetSummaryStats http.Handler
	GetDetailedLogs http.Handler
	CORS            http.Handler
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

// New instantiates HTTP handlers for all the HealthCheckChallenge service
// endpoints using the provided encoder and decoder. The handlers are mounted
// on the given mux using the HTTP verb and path defined in the design.
// errhandler is called whenever a response fails to be encoded. formatter is
// used to format errors returned by the service methods prior to encoding.
// Both errhandler and formatter are optional and can be nil.
func New(
	e *healthcheckchallenge.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"GetSummaryStats", "GET", "/healthcheck_challenge/summary_stats"},
			{"GetDetailedLogs", "GET", "/healthcheck_challenge/detailed_logs"},
			{"CORS", "OPTIONS", "/healthcheck_challenge/summary_stats"},
			{"CORS", "OPTIONS", "/healthcheck_challenge/detailed_logs"},
		},
		GetSummaryStats: NewGetSummaryStatsHandler(e.GetSummaryStats, mux, decoder, encoder, errhandler, formatter),
		GetDetailedLogs: NewGetDetailedLogsHandler(e.GetDetailedLogs, mux, decoder, encoder, errhandler, formatter),
		CORS:            NewCORSHandler(),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "HealthCheckChallenge" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.GetSummaryStats = m(s.GetSummaryStats)
	s.GetDetailedLogs = m(s.GetDetailedLogs)
	s.CORS = m(s.CORS)
}

// MethodNames returns the methods served.
func (s *Server) MethodNames() []string { return healthcheckchallenge.MethodNames[:] }

// Mount configures the mux to serve the HealthCheckChallenge endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountGetSummaryStatsHandler(mux, h.GetSummaryStats)
	MountGetDetailedLogsHandler(mux, h.GetDetailedLogs)
	MountCORSHandler(mux, h.CORS)
}

// Mount configures the mux to serve the HealthCheckChallenge endpoints.
func (s *Server) Mount(mux goahttp.Muxer) {
	Mount(mux, s)
}

// MountGetSummaryStatsHandler configures the mux to serve the
// "HealthCheckChallenge" service "getSummaryStats" endpoint.
func MountGetSummaryStatsHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleHealthCheckChallengeOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/healthcheck_challenge/summary_stats", f)
}

// NewGetSummaryStatsHandler creates a HTTP handler which loads the HTTP
// request and calls the "HealthCheckChallenge" service "getSummaryStats"
// endpoint.
func NewGetSummaryStatsHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeGetSummaryStatsRequest(mux, decoder)
		encodeResponse = EncodeGetSummaryStatsResponse(encoder)
		encodeError    = EncodeGetSummaryStatsError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "getSummaryStats")
		ctx = context.WithValue(ctx, goa.ServiceKey, "HealthCheckChallenge")
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

// MountGetDetailedLogsHandler configures the mux to serve the
// "HealthCheckChallenge" service "getDetailedLogs" endpoint.
func MountGetDetailedLogsHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleHealthCheckChallengeOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/healthcheck_challenge/detailed_logs", f)
}

// NewGetDetailedLogsHandler creates a HTTP handler which loads the HTTP
// request and calls the "HealthCheckChallenge" service "getDetailedLogs"
// endpoint.
func NewGetDetailedLogsHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeGetDetailedLogsRequest(mux, decoder)
		encodeResponse = EncodeGetDetailedLogsResponse(encoder)
		encodeError    = EncodeGetDetailedLogsError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "getDetailedLogs")
		ctx = context.WithValue(ctx, goa.ServiceKey, "HealthCheckChallenge")
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
// service HealthCheckChallenge.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleHealthCheckChallengeOrigin(h)
	mux.Handle("OPTIONS", "/healthcheck_challenge/summary_stats", h.ServeHTTP)
	mux.Handle("OPTIONS", "/healthcheck_challenge/detailed_logs", h.ServeHTTP)
}

// NewCORSHandler creates a HTTP handler which returns a simple 204 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})
}

// HandleHealthCheckChallengeOrigin applies the CORS response headers
// corresponding to the origin for the service HealthCheckChallenge.
func HandleHealthCheckChallengeOrigin(h http.Handler) http.Handler {
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
				w.WriteHeader(204)
				return
			}
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
