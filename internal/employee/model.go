package employee

import (
	"database/sql"
	"fmt"
	"strings"
)

// Employee struct TODO: Обсудить почему теги `json:"updated_at"` = зло?
type Employee struct {
	Id        string
	UserName  string
	FirstName string
	LastName  string
	CreatedAt string
	UpdatedAt string
}

var employeeColumns = []string{"id", "username", "first_name", "last_name", "created_at", "updated_at"}

func GetEmployeeByID(db *sql.DB, id int) (*Employee, error) {
	const op = "model.employee.GetEmployeeByID"
	var employee *Employee
	var q string
	if id <= 0 {
		return employee, fmt.Errorf("%s:%w", op, "Parametr id is empty")
	}

	q = fmt.Sprintf("SELECT %s FROM employee", strings.Join(employeeColumns, ","))

	row := db.QueryRow(q)

	employeeRow, err := employeeFromRows(row)
	if err != nil {
		return employee, fmt.Errorf("%s:%w", op, err)
	}

	return employeeRow, nil
}

func employeeFromRows(row *sql.Row) (*Employee, error) {
	var t Employee
	err := row.Scan(&t.Id, &t.UserName, &t.FirstName, &t.LastName, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}
