package app

import (
	"context"

	"github.com/maitesin/hermes/internal/domain"
)

type DeliveriesRepository interface {
	Insert(ctx context.Context, delivery domain.Delivery) error
	FindByTrackingID(ctx context.Context, trackingID string) (domain.Delivery, error)
	FindAllNotDelivered(ctx context.Context) ([]domain.Delivery, error)
}
