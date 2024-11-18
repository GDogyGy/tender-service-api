package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"TenderServiceApi/internal/model"
)

type TenderRepository struct {
	db *sql.DB
}

func NewTenderRepository(db *sql.DB) *TenderRepository {
	return &TenderRepository{db: db}
}

func (t *TenderRepository) GetTenderList(serviceType string) ([]model.Tender, error) {
	const op = "service.tender.GetTenderList"
	var tenders []model.Tender
	var rows *sql.Rows
	var err error

	if len(serviceType) > 0 {
		st := strings.Split(serviceType, ",")
		//rows, err = t.db.Query(`SELECT * FROM tender WHERE tender.service_type IN ($1)`, strings.Join(st, "', '")) // TODO: так почему то не работает: параметр не подхватывает возможно дело в github.com/lib/pq
		rows, err = t.db.Query(fmt.Sprintf("SELECT %s FROM tender WHERE tender.service_type IN ('%s')", strings.Join(teC, ","), strings.Join(st, "', '")))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return tenders, fmt.Errorf("%s:%w", op, err)
			}
			return tenders, err
		}
	} else {
		rows, err = t.db.Query("SELECT * FROM tender")
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return tenders, fmt.Errorf("%s:%w", op, err)
			}
			return tenders, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		tender, err := t.tenderFromRows(rows)
		if err != nil {
			return tenders, fmt.Errorf("%s:%w", op, err)
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

var teC = []string{"id", "name", "description", "service_type", "status", "responsible"}

func (t *TenderRepository) tenderFromRows(row *sql.Rows) (model.Tender, error) {
	var te model.Tender
	err := row.Scan(&te.Id, &te.Name, &te.Description, &te.ServiceType, &te.Status, &te.Responsible)
	return te, err
}
