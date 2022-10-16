package app

import (
	"log"
	"strings"

	"github.com/maitesin/hermes/pkg/comm"
	"github.com/maitesin/hermes/pkg/tracker"
)

func Listen(t tracker.Tracker) comm.Handler {
	return func(message comm.Message) error {
		events, err := t.Track(strings.TrimSpace(message.Text))
		if err != nil {
			return err
		}

		log.Printf("[events] %#v", events)

		return nil
	}
}
