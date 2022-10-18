package app

import "strings"

type Delivery struct {
	TrackingID     string
	Log            string
	ConversationID int64
	Delivered      bool
}

func NewDelivery(trackingID, log string, conversationID int64, delivered bool) Delivery {
	return Delivery{
		TrackingID:     trackingID,
		Log:            log,
		ConversationID: conversationID,
		Delivered:      delivered,
	}
}

func (d *Delivery) updateDelivered() {
	d.Delivered = strings.Contains(d.Log, "Entregado")
}
