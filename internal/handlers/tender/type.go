package tender

type tenderDTO struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"service_type"`
	Status      string `json:"status"`
	Responsible string `json:"responsible"`
}

// TODO валидно ли тут и так хранить?
type argCreatTender struct {
	Username       string `json:"creatorUsername"`
	OrganizationId string `json:"organizationId"`
}
