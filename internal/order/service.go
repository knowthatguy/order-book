package order

import (
	"context"
	"errors"
	"math"
	"reflect"
	"sort"
	"strings"

	cmap "github.com/orcaman/concurrent-map"
)

var (
	ErrOrderNotFound    = errors.New("Order not found")
	ErrOrderBookIsEmpty = errors.New("Order Book is empty")
)

// OrderService describes the Order service.
type OrderService interface {
	Create(ctx context.Context, order Order) (string, error)
	GetByID(ctx context.Context, id string) (Order, error)
	RemoveByID(ctx context.Context, id string) error
	GetMarketSnapshot(ctx context.Context) (MarketSnapshot, error)
}

// Order represents an order
type Order struct {
	ID        string  `json:"id,omitempty"`
	UUID      string  `json:"uuid,omitempty"`
	Side      string  `json:"side"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Status    string  `json:"status,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

// MarketSnapshot represents Market Data snapshot
type MarketSnapshot struct {
	Asks   []MarketSnapshotItem `json:"asks,omitempty"`
	Bids   []MarketSnapshotItem `json:"bids,omitempty"`
	Spread float64              `json:"spread"`
}

// MarketSnapshotItem represents item for Market Data snapshot
type MarketSnapshotItem struct {
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
}

// Sort the Asks and Bids orders
func (m MarketSnapshot) Sort() {
	sort.SliceStable(m.Asks, func(i, j int) bool {
		return m.Asks[i].Price <= m.Asks[j].Price
	})

	sort.SliceStable(m.Bids, func(i, j int) bool {
		return m.Bids[i].Price >= m.Bids[j].Price
	})
}

// Orders represent Orders
var Orders = cmap.New()

// Prices represents the min and max
type Prices struct {
	AskMinPrice, BidMaxPrice float64
	AskMinID, BidMaxID       string
}

func (p *Prices) setPrices() {
	var a, b float64
	var i, j string
	for order := range Orders.IterBuffered() {
		val := reflect.ValueOf(order.Val)
		if val.FieldByName("Status").String() == "Active" {
			if strings.ToUpper(val.FieldByName("Side").String()) == "BID" {
				b = math.Max(b, val.FieldByName("Price").Float())
				j = reflect.ValueOf(order.Key).String()
			} else {
				a = math.Min(a, val.FieldByName("Price").Float())
				i = reflect.ValueOf(order.Key).String()
			}
		}
	}

	(*p).BidMaxPrice = b
	(*p).BidMaxID = j
	(*p).AskMinPrice = a
	(*p).AskMinID = i
}

func (p Prices) getSpread() float64 {
	return p.AskMinPrice - p.BidMaxPrice
}

var spread Prices
