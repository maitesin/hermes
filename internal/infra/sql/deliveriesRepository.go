package sql

import (
	"context"

	"github.com/maitesin/hermes/internal/app"
	"github.com/maitesin/hermes/internal/domain"
	"github.com/upper/db/v4"
)

const (
	deliveryTable = "deliveries"
)

type Delivery struct {
	TrackingID string `db:"tracking_id"`
	Log        string `db:"log"`
	Delivered  bool   `db:"delivered"`
}

type DeliveriesRepository struct {
	sess db.Session
}

func NewDeliveriesRepository(sess db.Session) *DeliveriesRepository {
	return &DeliveriesRepository{sess: sess}
}

func (dr *DeliveriesRepository) Insert(ctx context.Context, delivery domain.Delivery) error {
	sDelivery := domain2SQLDelivery(delivery)

	_, err := dr.sess.WithContext(ctx).
		Collection(deliveryTable).
		Insert(sDelivery)

	return err
}

func (dr *DeliveriesRepository) FindByTrackingID(ctx context.Context, trackingID string) (domain.Delivery, error) {
	var delivery Delivery
	err := dr.sess.WithContext(ctx).
		Collection(deliveryTable).
		Find(db.Cond{"tracking_id": trackingID}).
		One(&delivery)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return domain.Delivery{}, app.NewDeliveryNotFoundError(trackingID)
		}
		return domain.Delivery{}, err
	}
	return domain.NewDelivery(delivery.TrackingID, delivery.Log), nil
}

func (dr *DeliveriesRepository) FindAllNotDelivered(ctx context.Context) ([]domain.Delivery, error) {
	var deliveries []Delivery
	err := dr.sess.WithContext(ctx).
		Collection(deliveryTable).
		Find(db.Cond{"delivered": false}).
		All(&deliveries)
	if err != nil {
		return nil, err
	}
	return sql2DomainDeliveries(deliveries...), nil
}

func sql2DomainDeliveries(sDeliveries ...Delivery) []domain.Delivery {
	deliveries := make([]domain.Delivery, len(sDeliveries))

	for i := range sDeliveries {
		deliveries[i] = domain.NewDelivery(sDeliveries[i].TrackingID, sDeliveries[i].Log)
	}

	return deliveries
}

func domain2SQLDelivery(dDelivery domain.Delivery) Delivery {
	return Delivery{
		TrackingID: dDelivery.TrackingID,
		Log:        dDelivery.Log,
	}
}
