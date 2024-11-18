package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"TenderServiceApi/internal/model"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (e *EmployeeRepository) GetEmployeeByID(id string) (model.Employee, error) {
	const op = "model.employee.GetEmployeeByID"
	var employee model.Employee
	if len(id) <= 0 {
		return employee, fmt.Errorf("%s:%w", op, "Parametr id is empty")
	}
	// TODO: проблема с передачей аргументов чистой функции QueryRow из за pq
	row := e.db.QueryRow(fmt.Sprintf("SELECT \"id\", \"username\", \"first_name\", \"last_name\", \"created_at\", \"updated_at\" FROM employee WHERE id ='%s'", id))

	employeeRow, err := e.employeeFromRows(row)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return employee, fmt.Errorf("%s:%w", op, err)
		}
		return employee, err
	}

	return employeeRow, nil
}

var emC = []string{"id", "username", "first_name", "last_name", "created_at", "updated_at"}

func (e *EmployeeRepository) employeeFromRows(row *sql.Row) (model.Employee, error) {
	var em model.Employee
	err := row.Scan(&em.Id, &em.UserName, &em.FirstName, &em.LastName, &em.CreatedAt, &em.UpdatedAt)
	return em, err
}
