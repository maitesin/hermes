package app

import (
	"fmt"
	"net/http"
)

func Asdf(repository DeliveriesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveries, err := repository.FindAllNotDelivered(r.Context())
		if err != nil {
			// TODO: return nice HTTP error
			panic(err)
		}
		w.Write([]byte("<h1>Undelivered</h1>"))
		for _, delivery := range deliveries {
			w.Write([]byte(fmt.Sprintf("<h2>%s</h2><br>%s", delivery.TrackingID, delivery.Log)))
		}
	}
}
