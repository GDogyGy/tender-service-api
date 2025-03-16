package model

type Bids struct {
	Id          string
	Name        string
	Description string
	Status      string
	TenderId    string
	Version     int
	Responsible string
}

func (a *Bids) FillDefault(bids Bids) {
	if a.Id == "" {
		a.Id = bids.Id
	}

	if a.Name == "" {
		a.Name = bids.Name
	}

	if a.Description == "" {
		a.Description = bids.Description
	}

	if a.Status == "" {
		a.Status = bids.Status
	}

	if a.TenderId == "" {
		a.TenderId = bids.TenderId
	}

	if a.Version == 0 {
		a.Version = bids.Version
	}

	if a.Responsible == "" {
		a.Responsible = bids.Responsible
	}
}
