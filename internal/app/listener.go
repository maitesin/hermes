package app

import (
	"github.com/maitesin/hermes/pkg/comm"
	"github.com/maitesin/hermes/pkg/tracker"
)

func Listen([]tracker.Tracker) func(handler comm.Handler) error {
	return func(handler comm.Handler) error {
		return nil
	}
}
