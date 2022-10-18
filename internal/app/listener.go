package app

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/maitesin/hermes/pkg/comm"
	"github.com/maitesin/hermes/pkg/tracker"
)

func Listen(ctx context.Context, t tracker.Tracker, repo DeliveriesRepository) comm.Handler {
	return func(message comm.Message) error {
		trackingID := strings.TrimSpace(message.Text)
		events, err := t.Track(trackingID)
		if err != nil {
			return err
		}

		eventsLog := fmt.Sprintf("%v", events)
		log.Printf("[events] %s", eventsLog)

		delivery := NewDelivery(trackingID, eventsLog, message.Conversation, false)
		delivery.updateDelivered()

		err = repo.Insert(ctx, delivery)
		if err != nil {
			return err
		}

		return nil
	}
}
