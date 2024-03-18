package app

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/maitesin/hermes/pkg/tracker/group"

	"github.com/maitesin/hermes/pkg/comm"
)

func Listen(ctx context.Context, g group.Group, repo DeliveriesRepository) comm.Handler {
	return func(message comm.Message) error {
		trackingID := strings.TrimSpace(message.Text)
		name, events, delivered, err := g.Track(trackingID)
		if err != nil {
			return err
		}

		eventsLog := fmt.Sprintf("%v", events)
		log.Printf("[events] %s", eventsLog)

		err = repo.Insert(ctx, NewDelivery(name, trackingID, eventsLog, message.Conversation, delivered))
		if err != nil {
			return err
		}

		return nil
	}
}
