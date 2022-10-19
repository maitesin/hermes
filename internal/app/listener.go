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
		events, delivered, err := t.Track(trackingID)
		if err != nil {
			return err
		}

		eventsLog := fmt.Sprintf("%v", events)
		log.Printf("[events] %s", eventsLog)

		err = repo.Insert(ctx, NewDelivery(trackingID, eventsLog, message.Conversation, delivered))
		if err != nil {
			return err
		}

		return nil
	}
}
