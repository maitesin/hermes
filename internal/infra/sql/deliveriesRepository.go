package sql

import (
	"context"

	"github.com/maitesin/hermes/internal/app"
	"github.com/upper/db/v4"
)

const (
	deliveryTable = "deliveries"
)

func onConflictUpdateLogAndDelivered(queryIn string) string {
	return queryIn + `ON CONFLICT (tracking_id) DO UPDATE SET
        log = EXCLUDED.log,
        delivered = EXCLUDED.delivered,
		conversation_id = EXCLUDED.conversation_id,
		courier = EXCLUDED.courier
`
}

type Delivery struct {
	TrackingID     string `db:"tracking_id"`
	Log            string `db:"log"`
	ConversationID int64  `db:"conversation_id"`
	Delivered      bool   `db:"delivered"`
	Courier        string `db:"courier"`
}

type DeliveriesRepository struct {
	sess db.Session
}

func NewDeliveriesRepository(sess db.Session) *DeliveriesRepository {
	return &DeliveriesRepository{sess: sess}
}

func (dr *DeliveriesRepository) Insert(ctx context.Context, delivery app.Delivery) error {
	sDelivery := app2SQLDelivery(delivery)

	_, err := dr.sess.WithContext(ctx).
		SQL().
		InsertInto(deliveryTable).
		Values(sDelivery).
		Amend(onConflictUpdateLogAndDelivered).
		Exec()

	return err
}

func (dr *DeliveriesRepository) FindByTrackingID(ctx context.Context, trackingID string) (app.Delivery, error) {
	var delivery Delivery
	err := dr.sess.WithContext(ctx).
		Collection(deliveryTable).
		Find(db.Cond{"tracking_id": trackingID}).
		One(&delivery)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return app.Delivery{}, app.NewDeliveryNotFoundError(trackingID)
		}
		return app.Delivery{}, err
	}
	return app.NewDelivery(delivery.Courier, delivery.TrackingID, delivery.Log, delivery.ConversationID, delivery.Delivered), nil
}

func (dr *DeliveriesRepository) FindAllNotDelivered(ctx context.Context) ([]app.Delivery, error) {
	var deliveries []Delivery
	err := dr.sess.WithContext(ctx).
		Collection(deliveryTable).
		Find(db.Cond{"delivered": false}).
		All(&deliveries)
	if err != nil {
		return nil, err
	}
	return sql2AppDeliveries(deliveries...), nil
}

func sql2AppDeliveries(sDeliveries ...Delivery) []app.Delivery {
	deliveries := make([]app.Delivery, len(sDeliveries))

	for i := range sDeliveries {
		deliveries[i] = app.NewDelivery(
			sDeliveries[i].Courier,
			sDeliveries[i].TrackingID,
			sDeliveries[i].Log,
			sDeliveries[i].ConversationID,
			sDeliveries[i].Delivered,
		)
	}

	return deliveries
}

func app2SQLDelivery(dDelivery app.Delivery) Delivery {
	return Delivery{
		TrackingID:     dDelivery.TrackingID,
		Log:            dDelivery.Log,
		ConversationID: dDelivery.ConversationID,
		Delivered:      dDelivery.Delivered,
		Courier:        dDelivery.Courier,
	}
}
