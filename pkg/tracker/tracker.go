package tracker

type Tracker interface {
	Track(id string) ([]DeliveryEvent, bool, error)
}

type DeliveryEvent struct {
	Timestamp   string
	Information string
}
