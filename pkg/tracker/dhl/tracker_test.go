package dhl_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/maitesin/hermes/pkg/tracker/dhl"

	"github.com/maitesin/hermes/pkg/tracker"
	"github.com/stretchr/testify/require"
)

type roundTripperMock struct {
	Response *http.Response
	RespErr  error
}

func (rtm *roundTripperMock) RoundTrip(*http.Request) (*http.Response, error) {
	return rtm.Response, rtm.RespErr
}

func workingClient(reader io.Reader) *http.Client {
	return &http.Client{
		Transport: &roundTripperMock{
			Response: &http.Response{
				Body: io.NopCloser(reader),
			},
		},
	}
}

func failingClient() *http.Client {
	return &http.Client{
		Transport: &roundTripperMock{},
	}
}

type readerMock struct {
	n   int
	err error
}

func (rm *readerMock) Read([]byte) (n int, err error) {
	return rm.n, rm.err
}

func validUndeliveredPayload() string {
	return `{
  "shipments": [
    {
      "events": [
        {
          "timestamp": "2024-03-09T05:13:33+01:00",
          "statusCode": "transit",
          "status": "32",
          "description": "The shipment has left transit country",
          "pieceIds": [
            "JJATA8235789000230623"
          ]
        },
        {
          "timestamp": "2024-03-08T02:24:33Z",
          "statusCode": "transit",
          "status": "7",
          "description": "The shipment has been dropped off at sorting center by sender",
          "pieceIds": [
            "JJATA8235789000230623"
          ]
        }
      ]
    }
  ]
}`
}

func validDeliveredPayload() string {
	return `{
  "shipments": [
    {
      "events": [
        {
          "timestamp": "2024-03-09T05:13:33+01:00",
          "statusCode": "transit",
          "status": "32",
          "description": "The shipment has left transit country",
          "pieceIds": [
            "JJATA8235789000230623"
          ]
        },
        {
          "timestamp": "2024-03-08T02:24:33Z",
          "statusCode": "transit",
          "status": "7",
          "description": "Delivered",
          "pieceIds": [
            "JJATA8235789000230623"
          ]
        }
      ]
    }
  ]
}`
}

func tooManyShipmentsPayload() string {
	return `{
  "shipment": [
    {
      "events": [
        {
          "timestamp": "2024-03-09T05:13:33+01:00",
          "statusCode": "transit",
          "status": "32",
          "description": "The shipment has left transit country",
          "pieceIds": [
            "JJATA8235789000230623"
          ]
        }
      ]
    },
    {
      "events": [
        {
          "timestamp": "2024-03-08T02:24:33Z",
          "statusCode": "transit",
          "status": "7",
          "description": "Delivered",
          "pieceIds": [
            "JJATA8235789000230623"
          ]
        }
      ]
    }
  ]
}`
}

func TestTracker_Track(t *testing.T) {
	tests := []struct {
		name          string
		client        *http.Client
		wantEvents    []tracker.DeliveryEvent
		wantDelivered bool
		wantErr       error
	}{
		{
			name: `Given a working HTTP client and a valid tracker ID,
                   when the method Track is called,
                   then it returns a list of delivery events that does not contain the delivery event`,
			client: workingClient(bytes.NewBufferString(validUndeliveredPayload())),
			wantEvents: []tracker.DeliveryEvent{
				{
					Timestamp:   "2024-03-09T05:13:33+01:00",
					Information: "The shipment has left transit country",
				},
				{
					Timestamp:   "2024-03-08T02:24:33Z",
					Information: "The shipment has been dropped off at sorting center by sender",
				},
			},
			wantDelivered: false,
			wantErr:       nil,
		},
		{
			name: `Given a working HTTP client and a valid tracker ID,
                   when the method Track is called,
                   then it returns a list of delivery events that contains the delivery event`,
			client: workingClient(bytes.NewBufferString(validDeliveredPayload())),
			wantEvents: []tracker.DeliveryEvent{
				{
					Timestamp:   "2024-03-09T05:13:33+01:00",
					Information: "The shipment has left transit country",
				},
				{
					Timestamp:   "2024-03-08T02:24:33Z",
					Information: "Delivered",
				},
			},
			wantDelivered: true,
			wantErr:       nil,
		},
		{
			name: `Given a failing HTTP client and a valid tracker ID,
                   when the method Track is called,
                   then it returns an error regarding the HTTP client not working`,
			client:        failingClient(),
			wantEvents:    nil,
			wantDelivered: false,
			wantErr:       &url.Error{},
		},
		{
			name: `Given a failing HTTP client and a valid tracker ID,
                   when the method Track is called,
                   then it returns an error regarding the HTTP client not working`,
			client:        workingClient(&readerMock{n: 0, err: errors.New("")}),
			wantEvents:    nil,
			wantDelivered: false,
			wantErr:       errors.New(""),
		},
		{
			name: `Given a working HTTP client and a valid tracker ID,
                   when the method Track is called and the response is an invalid JSON,
                   then it returns an error regarding the response body`,
			client:        workingClient(bytes.NewBufferString("")),
			wantEvents:    nil,
			wantDelivered: false,
			wantErr:       &json.SyntaxError{},
		},
		{
			name: `Given a working HTTP client and a valid tracker ID,
                   when the method Track is called and the response contains multiple shipment information,
                   then it returns an error regarding the response containing multiple shipment information`,
			client:        workingClient(bytes.NewBufferString(tooManyShipmentsPayload())),
			wantEvents:    nil,
			wantDelivered: false,
			wantErr:       errors.New(""),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trackerTest, err := dhl.NewTracker(tt.client, "")
			require.Nil(t, err)

			gotEvents, gotDelivered, err := trackerTest.Track("does not matter regarding testing")
			if tt.wantErr != nil {
				require.ErrorAs(t, err, &tt.wantErr)
			} else {
				require.Nil(t, err)
			}

			require.Equal(t, tt.wantEvents, gotEvents)
			require.Equal(t, tt.wantDelivered, gotDelivered)
		})
	}
}
