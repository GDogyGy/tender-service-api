package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"TenderServiceApi/internal/model"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (e *Repository) FetchById(ctx context.Context, id string) (model.Employee, error) {
	const op = "repository.employee.FetchById"
	var employee model.Employee
	if len(id) <= 0 {
		return employee, fmt.Errorf("%s:%s", op, "Parametr id is empty")
	}

	r := e.db.QueryRowx(fmt.Sprintf(`SELECT %s FROM employee WHERE id = $1`, strings.Join(column, ", ")), id)

	employeeRow, err := e.fromRow(r)

	if errors.Is(err, sql.ErrNoRows) {
		return employee, fmt.Errorf("%s:%w", op, err)
	}

	if err != nil {
		return employee, err
	}

	return employeeRow, nil
}

func (e *Repository) FetchByUserName(ctx context.Context, username string) (model.Employee, error) {
	const op = "repository.employee.FetchByUserName"
	var employee model.Employee
	if len(username) <= 0 {
		return employee, fmt.Errorf("%s:%s", op, "Parametr id is empty")
	}

	r := e.db.QueryRowx(fmt.Sprintf(`SELECT %s FROM employee WHERE username = $1`, strings.Join(column, ", ")), username)

	employeeRow, err := e.fromRow(r)

	if errors.Is(err, sql.ErrNoRows) {
		return employee, fmt.Errorf("%s:%w", op, err)
	}

	if err != nil {
		return employee, err
	}

	return employeeRow, nil
}

var column = []string{"id", "username", "first_name", "last_name", "created_at", "updated_at"}

func (e *Repository) fromRows(rows *sqlx.Rows) (model.Employee, error) {
	var r row
	err := rows.StructScan(&r)
	m := r.toModel()
	return m, err
}

func (e *Repository) fromRow(ro *sqlx.Row) (model.Employee, error) {
	var r row
	err := ro.StructScan(&r)
	m := r.toModel()
	return m, err
}
