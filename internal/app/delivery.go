package app

type Delivery struct {
	Courier        string
	TrackingID     string
	Log            string
	ConversationID int64
	Delivered      bool
}

func NewDelivery(courier, trackingID, log string, conversationID int64, delivered bool) Delivery {
	return Delivery{
		Courier:        courier,
		TrackingID:     trackingID,
		Log:            log,
		ConversationID: conversationID,
		Delivered:      delivered,
	}
}
