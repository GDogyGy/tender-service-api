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

type bidsFeedbackDTO struct {
	Id          string `json:"id"`
	BidID       string `json:"bidID"`
	Description string `json:"description"`
	Responsible string `json:"responsible"`
	CreatedAt   string `json:"created_at"`
}

type argCreatBids struct {
	Username       string `json:"creatorUsername"`
	OrganizationId string `json:"organizationId"`
}
