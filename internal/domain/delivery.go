package domain

type Delivery struct {
	TrackingID string
	Log        string
}

func NewDelivery(trackingID, log string) Delivery {
	return Delivery{
		TrackingID: trackingID,
		Log:        log,
	}
}
