package order

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
)

func Test_service_Create(t *testing.T) {
	type fields struct {
		logger log.Logger
	}
	type args struct {
		ctx      context.Context
		newOrder Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				logger: tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.newOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("service.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
