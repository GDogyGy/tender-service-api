package update

type bidDTO struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	TenderId    string `json:"tender_id"`
	Version     int    `json:"version"`
	Responsible string `json:"responsible"`
}
