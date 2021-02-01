package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	api "github.com/knowthatguy/order-book/internal/api"
	order "github.com/knowthatguy/order-book/internal/order"
)

var (
	ErrBadRouting = errors.New("bad routing")
)

// NewHTTPService wires Go kit endpoints to the HTTP transport.
func NewHTTPService(
	svcEndpoints api.Endpoints, options []kithttp.ServerOption, logger log.Logger,
) http.Handler {
	// set-up router and initialize http endpoints
	var (
		r            = mux.NewRouter()
		errorLogger  = kithttp.ServerErrorLogger(logger)
		errorEncoder = kithttp.ServerErrorEncoder(encodeErrorResponse)
	)
	options = append(options, errorLogger, errorEncoder)

	// HTTP Get - /orders
	r.Methods("GET").Path("/orders").Handler(kithttp.NewServer(
		svcEndpoints.GetMarketSnapshot,
		decodeGetMarketSnapshotRequest,
		encodeResponse,
		options...,
	))

	// HTTP Post - /orders
	r.Methods("POST").Path("/orders").Handler(kithttp.NewServer(
		svcEndpoints.Create,
		decodeCreateRequest,
		encodeResponse,
		options...,
	))

	// HTTP Get - /orders/{id}
	r.Methods("GET").Path("/orders/{id}").Handler(kithttp.NewServer(
		svcEndpoints.GetByID,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	))

	// HTTP Post - /orders/{id}
	r.Methods("DELETE").Path("/orders/{id}").Handler(kithttp.NewServer(
		svcEndpoints.RemoveByID,
		decodeRemoveByIDRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeGetMarketSnapshotRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req api.GetMarketSnapshotRequest
	return req, nil
}

func decodeCreateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req api.CreateRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Order); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return api.GetByIDRequest{ID: id}, nil
}

func decodeRemoveByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return api.RemoveByIDRequest{ID: id}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case order.ErrOrderNotFound:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
