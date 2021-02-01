package api

import order "github.com/knowthatguy/order-book/internal/order"

// CreateRequest holds the request parameters for the Create method.
type CreateRequest struct {
	Order order.Order
}

// CreateResponse holds the response values for the Create method.
type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

// GetByIDRequest holds the request parameters for the GetByID method.
type GetByIDRequest struct {
	ID string
}

// GetByIDResponse holds the response values for the GetByID method.
type GetByIDResponse struct {
	Order order.Order `json:"order"`
	Err   error       `json:"error,omitempty"`
}

// RemoveByIDRequest holds the request parameters for the RemoveByID method.
type RemoveByIDRequest struct {
	ID string
}

// RemoveByIDResponse holds the response values for the RemoveByID method.
type RemoveByIDResponse struct {
	Err error `json:"error,omitempty"`
}

// GetMarketSnapshotRequest holds
type GetMarketSnapshotRequest struct{}

// GetMarketSnapshotResponse holds the response values for the GetByID method.
type GetMarketSnapshotResponse struct {
	MarketSnapshot order.MarketSnapshot `json:"result"`
	Err            error                `json:"error,omitempty"`
}
