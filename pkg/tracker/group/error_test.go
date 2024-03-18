package group_test

import (
	"fmt"
	"testing"

	"github.com/maitesin/hermes/pkg/tracker/group"
	"github.com/stretchr/testify/require"
)

func TestDeliveryNotFoundError_Error(t *testing.T) {
	t.Parallel()

	trackingID := "12345678"
	err := group.NewTrackerNotFoundForDelivery(trackingID)

	require.ErrorAs(t, err, &group.TrackerNotFoundForDelivery{})
	require.Equal(t, fmt.Sprintf("tracker not found for delivery with tracking ID %q", trackingID), err.Error())
}
