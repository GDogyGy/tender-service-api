package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type Employee struct {
	Id        string `json:"id"`
	UserName  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var employeeColumns = []string{"id", "username", "first_name", "last_name", "created_at", "updated_at"}

func employeeFromRows(row *sql.Row) (*Employee, error) {
	var t Employee
	err := row.Scan(&t.Id, &t.UserName, &t.FirstName, &t.LastName, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}

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
