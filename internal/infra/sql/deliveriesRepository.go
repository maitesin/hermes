package sql

import (
	"context"
	"database/sql"

	"github.com/maitesin/hermes/internal/app"
)

type Delivery struct {
	TrackingID     string `db:"tracking_id"`
	Log            string `db:"log"`
	ConversationID int64  `db:"conversation_id"`
	Delivered      bool   `db:"delivered"`
	Courier        string `db:"courier"`
}

type DeliveriesRepository struct {
	db *sql.DB
}

func NewDeliveriesRepository(db *sql.DB) *DeliveriesRepository {
	return &DeliveriesRepository{db: db}
}

func (dr *DeliveriesRepository) Insert(ctx context.Context, delivery app.Delivery) error {
	sDelivery := app2SQLDelivery(delivery)

	_, err := dr.db.ExecContext(
		ctx,
		"INSERT INTO deliveries (tracking_id, log, conversation_id, courier, delivered) VALUES (?, ?, ?, ?, ?)",
		sDelivery.TrackingID,
		sDelivery.Log,
		sDelivery.ConversationID,
		sDelivery.Courier,
		sDelivery.Delivered,
	)

	return err
}

func (dr *DeliveriesRepository) Update(ctx context.Context, delivery app.Delivery) error {
	sDelivery := app2SQLDelivery(delivery)

	_, err := dr.db.ExecContext(
		ctx,
		"UPDATE deliveries SET delivered = ?, log = ? WHERE tracking_id = ? and conversation_id = ?",
		sDelivery.Delivered,
		sDelivery.Log,
		sDelivery.TrackingID,
		sDelivery.ConversationID,
	)
	return err
}

func (dr *DeliveriesRepository) FindByTrackingID(ctx context.Context, trackingID string) (app.Delivery, error) {
	var delivery Delivery
	row, err := dr.db.QueryContext(ctx, "SELECT tracking_id, log, conversation_id, courier, delivered FROM deliveries WHERE tracking_id = ?", trackingID)
	if err != nil {
		return app.Delivery{}, err
	}
	defer row.Close()
	err = row.Scan(&delivery.TrackingID, &delivery.Log, &delivery.ConversationID, &delivery.Delivered)
	return app.NewDelivery(delivery.Courier, delivery.TrackingID, delivery.Log, delivery.ConversationID, delivery.Delivered), err
}

func (dr *DeliveriesRepository) FindAllNotDelivered(ctx context.Context) ([]app.Delivery, error) {
	var deliveries []Delivery
	rows, err := dr.db.QueryContext(ctx, "SELECT tracking_id, log, conversation_id, courier, delivered FROM deliveries WHERE delivered = FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var delivery Delivery
		err = rows.Scan(&delivery.TrackingID, &delivery.Log, &delivery.ConversationID, &delivery.Courier, &delivery.Delivered)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, delivery)
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
