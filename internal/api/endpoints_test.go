package api

import (
	"reflect"
	"testing"

	order "github.com/knowthatguy/order-book/internal/order"
)

func TestMakeEndpoints(t *testing.T) {
	type args struct {
		s order.OrderService
	}
	tests := []struct {
		name string
		args args
		want Endpoints
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeEndpoints(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeEndpoints() = %v, want %v", got, tt.want)
			}
		})
	}
}
