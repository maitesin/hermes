package tracker

type Tracker interface {
	Track(id string) ([]DeliveryEvent, error)
}

type DeliveryEvent struct {
	Timestamp   string
	Information string
}
