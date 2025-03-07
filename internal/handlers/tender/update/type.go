package update

type tenderDTO struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"service_type"`
	Status      string `json:"status"`
	Version     int    `json:"version"`
	Responsible string `json:"responsible"`
}
