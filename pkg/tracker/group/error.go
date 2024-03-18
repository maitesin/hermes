package group

import "fmt"

type TrackerNotFoundForDelivery struct {
	trackingID string
}

func (t TrackerNotFoundForDelivery) Error() string {
	return fmt.Sprintf("tracker not found for delivery with tracking ID %q", t.trackingID)
}

func NewTrackerNotFoundForDelivery(trackingID string) TrackerNotFoundForDelivery {
	return TrackerNotFoundForDelivery{
		trackingID: trackingID,
	}
}
