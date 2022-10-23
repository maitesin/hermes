package app_test

import (
	"fmt"
	"testing"

	"github.com/maitesin/hermes/internal/app"
	"github.com/stretchr/testify/require"
)

func TestDeliveryNotFoundError_Error(t *testing.T) {
	t.Parallel()

	trackingID := "12345678"
	err := app.NewDeliveryNotFoundError(trackingID)

	require.ErrorAs(t, err, &app.DeliveryNotFoundError{})
	require.Equal(t, fmt.Sprintf("delivery with tracking ID %q not found", trackingID), err.Error())
}
