package employee

import (
	"TenderServiceApi/internal/model"
)

type row struct {
	Id        string `db:"id"`
	UserName  string `db:"username"`
	FirstName string `db:"first_name"`
	LastMame  string `db:"last_name"`
}

func (r *row) toModel() model.Employee {
	return model.Employee{
		Id:        r.Id,
		UserName:  r.UserName,
		FirstName: r.FirstName,
		LastName:  r.LastMame,
	}
}
