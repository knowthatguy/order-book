package order

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
)

// service implements the Order Service
type service struct {
	logger log.Logger
}

// NewService creates and returns a new Order service instance
func NewService(logger log.Logger) OrderService {
	return &service{
		logger: logger,
	}
}

func match(q int, p float64, s string) {
	var sideID string
	if strings.ToUpper(s) == "ASK" {
		sideID = spread.BidMaxID
	} else {
		sideID = spread.AskMinID
	}

	tmp, _ := Orders.Get(sideID)
	order := tmp.(Order)
	q = order.Quantity - q
	switch {
	case q < 0:
		fmt.Println(q)
		Orders.Remove(sideID)
		spread.setPrices()
		match(int(math.Abs(float64(q))), p, s)
	case q == 0:
		Orders.Remove(sideID)
		spread.setPrices()
	case q > 0:
		order.Quantity = q
		Orders.Remove(sideID)
		Orders.Set(sideID, order)
	}

}

// Create makes an order
func (s *service) Create(ctx context.Context, newOrder Order) (string, error) {
	uuid, _ := uuid.NewUUID()
	timestamp := time.Now().Unix()
	id := strconv.FormatInt(timestamp, 10)
	newOrder.ID = id
	newOrder.UUID = uuid.String()
	newOrder.Timestamp = timestamp
	side := newOrder.Side
	newOrder.Status = "Active"
	if !Orders.IsEmpty() {
		if strings.ToUpper(side) == "ASK" {
			// Matching
			if spread.BidMaxPrice >= newOrder.Price {
				match(newOrder.Quantity, newOrder.Price, side)
				newOrder.Status = "Completed"
			}
		} else {
			// Matching
			if spread.AskMinPrice <= newOrder.Price {
				match(newOrder.Quantity, newOrder.Price, side)
				newOrder.Status = "Completed"
			}
		}

	}

	Orders.Set(id, newOrder)
	spread.setPrices()
	return id, nil
}

// GetByID returns an order given by id
func (s *service) GetByID(ctx context.Context, id string) (Order, error) {
	logger := log.With(s.logger, "method", "GetByID")
	var order Order
	if tmp, ok := Orders.Get(id); ok {
		if !ok {
			level.Error(logger).Log("err", ok)
			return order, ErrOrderNotFound
		}
		order = tmp.(Order)
	}

	return order, nil
}

// RemoveByID changes the status of an order
func (s *service) RemoveByID(ctx context.Context, id string) error {
	logger := log.With(s.logger, "method", "RemovebyID")
	if Orders.IsEmpty() {
		level.Error(logger).Log("err", ErrOrderBookIsEmpty)
		return ErrOrderBookIsEmpty
	}

	if _, ok := Orders.Get(id); ok {
		if !ok {
			level.Error(logger).Log("err", ok)
			return ErrOrderNotFound
		}
	}
	Orders.Remove(id)
	spread.setPrices()
	return nil
}

// GetMarketSnapshot returns the snapshot with aggregated bids and asks
func (s *service) GetMarketSnapshot(ctx context.Context) (MarketSnapshot, error) {
	logger := log.With(s.logger, "method", "GetMarketSnapshot")
	snapshot := MarketSnapshot{}
	if Orders.IsEmpty() {
		level.Error(logger).Log("err", ErrOrderBookIsEmpty)
		return snapshot, ErrOrderBookIsEmpty
	}

	for order := range Orders.IterBuffered() {
		val := reflect.ValueOf(order.Val)

		new := MarketSnapshotItem{
			Price:    val.FieldByName("Price").Float(),
			Quantity: val.FieldByName("Quantity").Int(),
		}
		if val.FieldByName("Status").String() == "Active" {
			if strings.ToUpper(val.FieldByName("Side").String()) == "ASK" {
				snapshot.Asks = append(snapshot.Asks, new)
			} else {
				snapshot.Bids = append(snapshot.Bids, new)
			}
		}

	}

	// sorting
	snapshot.Sort()

	snapshot.Spread = spread.getSpread()
	return snapshot, nil
}
