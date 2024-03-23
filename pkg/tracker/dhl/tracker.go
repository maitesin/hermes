package dhl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/maitesin/hermes/pkg/tracker"
)

const urlRegex = "https://api-eu.dhl.com/track/shipments?trackingNumber=%s"

// Tracker for the DHL delivery service
type Tracker struct {
	client *http.Client
	key    string
}

// NewTracker constructor for the DHL tracker
func NewTracker(client *http.Client, key string) (*Tracker, error) {
	return &Tracker{
		client: client,
		key:    key,
	}, nil
}

func (t Tracker) Name() string {
	return "dhl"
}

func (t *Tracker) Track(id string) ([]tracker.DeliveryEvent, bool, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(urlRegex, id), nil)
	if err != nil {
		return nil, false, err
	}

	req.Header = http.Header{"DHL-API-Key": []string{t.key}}
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}

	var body body
	err = json.Unmarshal(respBody, &body)
	if err != nil {
		return nil, false, err
	}

	if len(body.Shipments) != 1 {
		return nil, false, fmt.Errorf(
			"expected information from a single shipment, found %d.\n%s",
			len(body.Shipments),
			body,
		)
	}

	var delivered bool
	events := make([]tracker.DeliveryEvent, len(body.Shipments[0].Events))
	for i, event := range body.Shipments[0].Events {
		events[i] = tracker.DeliveryEvent{
			Timestamp:   fmt.Sprintf("%s", event.Date),
			Information: fmt.Sprintf("%s", event.ExtendedText),
		}
		if !delivered {
			delivered = strings.Contains(events[i].Information, "elivered")
		}
	}

	return events, delivered, nil
}
