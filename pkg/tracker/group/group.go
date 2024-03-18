package group

import "github.com/maitesin/hermes/pkg/tracker"

type Group []tracker.Tracker

func NewGroup(trackers ...tracker.Tracker) Group {
	return trackers
}

func (g Group) Track(id string) (string, []tracker.DeliveryEvent, bool, error) {
	for _, t := range g {
		events, delivered, err := t.Track(id)
		if err != nil {
			continue
		}
		return t.Name(), events, delivered, err
	}
	return "", nil, false, NewTrackerNotFoundForDelivery(id)
}
