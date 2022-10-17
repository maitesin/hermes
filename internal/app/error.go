package app

import "fmt"

type DeliveryNotFoundError struct {
	trackingID string
}

func (d DeliveryNotFoundError) Error() string {
	return fmt.Sprintf("delivery with tracking ID %q not found", d.trackingID)
}

func NewDeliveryNotFoundError(trackingID string) DeliveryNotFoundError {
	return DeliveryNotFoundError{
		trackingID: trackingID,
	}
}
