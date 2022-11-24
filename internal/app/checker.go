package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/maitesin/hermes/pkg/comm"
	"github.com/maitesin/hermes/pkg/tracker"
)

func Checker(
	ctx context.Context,
	ticker *time.Ticker,
	t tracker.Tracker,
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
				events, delivered, err := t.Track(dbDelivery.TrackingID)
				if err != nil {
					return err
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
