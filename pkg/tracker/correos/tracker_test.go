package correos_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/maitesin/hermes/pkg/tracker"
	"github.com/maitesin/hermes/pkg/tracker/correos"
	"github.com/stretchr/testify/require"
)

type roundTripperMock struct {
	Response *http.Response
	RespErr  error
}

func (rtm *roundTripperMock) RoundTrip(*http.Request) (*http.Response, error) {
	return rtm.Response, rtm.RespErr
}

func clientResponding(response string) *http.Client {
	return &http.Client{
		Transport: &roundTripperMock{
			Response: &http.Response{
				Body: io.NopCloser(bytes.NewBufferString(response)),
			},
		},
	}
}

func validPayload() string {
	return `{
  "shipment": [
    {
      "events": [
        {
          "eventDate": "18/11/2021",
          "eventTime": "18:25:51",
          "phase": "2",
          "colour": "V",
          "summaryText": "Admitido",
          "extendedText": "El envío ha tenido admisión en origen",
          "actionWeb": null,
          "actionWebParam": null,
          "codired": "1254002"
        }
      ]
    }
  ]
}`
}

func TestTracker_Track(t *testing.T) {
	tests := []struct {
		name      string
		client    *http.Client
		trackerID string
		want      []tracker.DeliveryEvent
		wantErr   error
	}{
		{
			name:      "",
			client:    clientResponding(validPayload()),
			trackerID: "",
			want: []tracker.DeliveryEvent{
				{
					Timestamp:   "18/11/2021 18:25:51",
					Information: "Admitido (El envío ha tenido admisión en origen)",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			trackerTest, err := correos.NewTracker(tt.client)
			require.Nil(t, err)

			got, err := trackerTest.Track(tt.trackerID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.Nil(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
