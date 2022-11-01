package app

import (
	"context"
	"fmt"
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

				eventsLog := fmt.Sprintf("%v", events)
				if eventsLog == dbDelivery.Log {
					continue
				}

				msg := comm.Message{
					Conversation: dbDelivery.ConversationID,
					Text:         eventsLog,
				}
				//fmt.Printf("Message to be sent: %#v", msg)
				err = messenger.Message(
					msg,
				)
				if err != nil {
					return err
				}

				err = deliveriesRepository.Update(
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
