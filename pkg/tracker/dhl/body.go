package dhl

type body struct {
	Shipments []shipment `json:"shipments"`
}

type shipment struct {
	Events []event `json:"events"`
}

type event struct {
	Date         string `json:"timestamp"`
	ExtendedText string `json:"description"`
}
