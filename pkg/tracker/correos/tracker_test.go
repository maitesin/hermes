package correos_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
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

func workingClient(response string) *http.Client {
	return &http.Client{
		Transport: &roundTripperMock{
			Response: &http.Response{
				Body: io.NopCloser(bytes.NewBufferString(response)),
			},
		},
	}
}

func failingClient() *http.Client {
	return &http.Client{
		Transport: &roundTripperMock{},
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

func tooManyShipmentsPayload() string {
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
    },
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
		name    string
		client  *http.Client
		want    []tracker.DeliveryEvent
		wantErr error
	}{
		{
			name: `Given a working HTTP client and a valid tracker ID,
                   when the method Track is called,
                   then it returns a list of delivery events`,
			client: workingClient(validPayload()),
			want: []tracker.DeliveryEvent{
				{
					Timestamp:   "18/11/2021 18:25:51",
					Information: "Admitido (El envío ha tenido admisión en origen)",
				},
			},
			wantErr: nil,
		},
		{
			name: `Given a failing HTTP client and a valid tracker ID,
                   when the method Track is called,
                   then it returns an error regarding the HTTP client not working`,
			client:  failingClient(),
			want:    nil,
			wantErr: &url.Error{},
		},
		{
			name: `Given a working HTTP client and a valid tracker ID,
                   when the method Track is called and the response is an invalid JSON,
                   then it returns an error regarding the response body`,
			client:  workingClient(""),
			want:    nil,
			wantErr: &json.SyntaxError{},
		},
		{
			name: `Given a working HTTP client and a valid tracker ID,
                   when the method Track is called and the response contains multiple shipment information,
                   then it returns an error regarding the response containing multiple shipment information`,
			client:  workingClient(tooManyShipmentsPayload()),
			want:    nil,
			wantErr: errors.New(""),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trackerTest, err := correos.NewTracker(tt.client)
			require.Nil(t, err)

			got, err := trackerTest.Track("does not matter regarding testing")
			if tt.wantErr != nil {
				require.ErrorAs(t, err, &tt.wantErr)
			} else {
				require.Nil(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
