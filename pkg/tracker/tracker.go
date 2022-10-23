package tracker

//go:generate mockgen -destination=mocks/tracker.go -package=mocks . Tracker
type Tracker interface {
	Track(id string) ([]DeliveryEvent, bool, error)
}

type DeliveryEvent struct {
	Timestamp   string
	Information string
}
