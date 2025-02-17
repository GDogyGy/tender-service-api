package create

type bidsDTO struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	TenderId    string `json:"tenderId"`
	Version     int    `json:"version"`
	Responsible string `json:"responsible"`
}

type argCreatBids struct {
	Username       string `json:"creatorUsername"`
	OrganizationId string `json:"organizationId"`
}
