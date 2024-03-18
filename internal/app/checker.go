package app

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/maitesin/hermes/pkg/tracker/group"

	"github.com/maitesin/hermes/pkg/comm"
	"github.com/maitesin/hermes/pkg/tracker"
)

func Checker(
	ctx context.Context,
	ticker *time.Ticker,
	g group.Group,
	deliveriesRepository DeliveriesRepository,
	messenger comm.Messenger,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			deliveries, err := deliveriesRepository.FindAllNotDelivered(ctx)
			if err != nil {
				return err
			}

			for _, dbDelivery := range deliveries {
				name, events, delivered, err := g.Track(dbDelivery.TrackingID)
				if err != nil {
					log.Printf("Error found when getting information from trackers: %s", err.Error())
					continue
				}

				eventsLog := events2Log(events)
				if eventsLog != dbDelivery.Log {
					msg := comm.Message{
						Conversation: dbDelivery.ConversationID,
						Text:         fmt.Sprintf("%s:\n%s", dbDelivery.TrackingID, eventsLog),
					}

					err = messenger.Message(
						msg,
					)
					if err != nil {
						return err
					}
				}

				err = deliveriesRepository.Insert(
					ctx,
					NewDelivery(
						name,
						dbDelivery.TrackingID,
						eventsLog,
						dbDelivery.ConversationID,
						delivered,
					),
				)
				if err != nil {
					return err
				}
			}
		}
	}
}

func events2Log(events []tracker.DeliveryEvent) string {
	eventStrings := make([]string, len(events))

	for i, event := range events {
		eventStrings[i] = fmt.Sprintf("- %s:\n  %s\n", event.Timestamp, event.Information)
	}

	return strings.Join(eventStrings, "")
}
