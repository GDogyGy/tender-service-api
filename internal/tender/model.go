package tender

import (
	"TenderServiceApi/internal/employee"
)

type Tender struct {
	Id          string
	Name        string
	Description string
	Status      string
	ResponseId  string

	Employee *employee.Employee
}

var Columns = []string{"id", "name", "description", "service_type", "status", "responsible"}
