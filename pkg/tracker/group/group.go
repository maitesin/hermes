package group

import (
	"log"

	"github.com/maitesin/hermes/pkg/tracker"
)

type Group []tracker.Tracker

func NewGroup(trackers ...tracker.Tracker) Group {
	return trackers
}

func (g Group) Track(id string) (string, []tracker.DeliveryEvent, bool, error) {
	for _, t := range g {
		events, delivered, err := t.Track(id)
		if err != nil {
			log.Printf("Tracker %s couldn't find %s with error: %s", t.Name(), id, err.Error())
			continue
		}
		return t.Name(), events, delivered, err
	}
	return "", nil, false, NewTrackerNotFoundForDelivery(id)
}
