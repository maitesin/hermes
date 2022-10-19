package app

import (
	"context"
)

type DeliveriesRepository interface {
	Insert(ctx context.Context, delivery Delivery) error
	Update(ctx context.Context, delivery Delivery) error
	FindByTrackingID(ctx context.Context, trackingID string) (Delivery, error)
	FindAllNotDelivered(ctx context.Context) ([]Delivery, error)
}
