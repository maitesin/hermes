package correos

type body struct {
	Shipments []shipment `json:"shipment"`
}

type shipment struct {
	Events []event `json:"events"`
}

type event struct {
	Date         string `json:"eventDate"`
	Time         string `json:"eventTime"`
	SummaryText  string `json:"summaryText"`
	ExtendedText string `json:"extendedText"`
}
