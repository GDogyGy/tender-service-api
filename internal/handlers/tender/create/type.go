package create

type tenderDTO struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"service_type"`
	Status      string `json:"status"`
	Version     int    `json:"version"`
	Responsible string `json:"responsible"`
}

type argCreatTender struct {
	Username       string `json:"creatorUsername"`
	OrganizationId string `json:"organizationId"`
}
