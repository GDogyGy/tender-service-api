package model

type BidFeedback struct {
	Id          string
	BidID       string
	Description string
	Responsible string
	CreatedAt   string
}

func (a *BidFeedback) FillDefault(defaults *BidFeedback) {
	type field struct {
		src  interface{}
		dest interface{}
	}

	fields := []field{
		{&a.Id, &defaults.Id},
		{&a.BidID, &defaults.BidID},
		{&a.Description, &defaults.Description},
		{&a.Responsible, &defaults.Responsible},
		{&a.Responsible, &defaults.Responsible},
		{&a.CreatedAt, &defaults.CreatedAt},
	}

	for _, f := range fields {
		switch dest := f.dest.(type) {
		case *string:
			if *dest == "" {
				*dest = *f.src.(*string)
			}
		case *int:
			if *dest == 0 {
				*dest = *f.src.(*int)
			}
		}
	}
}
