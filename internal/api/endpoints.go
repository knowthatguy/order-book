package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	order "github.com/knowthatguy/order-book/internal/order"
)

// Endpoints holds all Go kit endpoints for the Order service.
type Endpoints struct {
	Create            endpoint.Endpoint
	GetByID           endpoint.Endpoint
	RemoveByID        endpoint.Endpoint
	GetMarketSnapshot endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the Order service.
func MakeEndpoints(s order.OrderService) Endpoints {
	return Endpoints{
		Create:            makeCreateEndpoint(s),
		GetByID:           makeGetByIDEndpoint(s),
		RemoveByID:        makeRemoveByIDEndpoint(s),
		GetMarketSnapshot: makeGetMarketSnapshotEndpoint(s),
	}
}

func makeCreateEndpoint(s order.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest) // type assertion
		id, err := s.Create(ctx, req.Order)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s order.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		orderRes, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Order: orderRes, Err: err}, err
	}
}

func makeRemoveByIDEndpoint(s order.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RemoveByIDRequest)
		err := s.RemoveByID(ctx, req.ID)
		return RemoveByIDResponse{Err: err}, nil
	}
}

func makeGetMarketSnapshotEndpoint(s order.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		snapshot, err := s.GetMarketSnapshot(ctx)
		return GetMarketSnapshotResponse{MarketSnapshot: snapshot, Err: err}, err
	}
}
