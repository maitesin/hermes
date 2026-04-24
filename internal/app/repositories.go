package app

import (
	"context"
)

//go:generate mockgen -destination=mocks/deliveries_repository.go -package=mocks . DeliveriesRepository

type DeliveriesRepository interface {
	Insert(ctx context.Context, delivery Delivery) error
	Update(ctx context.Context, delivery Delivery) error
	FindByTrackingID(ctx context.Context, trackingID string) (Delivery, error)
	FindAllNotDelivered(ctx context.Context) ([]Delivery, error)
}
