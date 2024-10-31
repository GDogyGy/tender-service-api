package tender

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"TenderServiceApi/internal/employee"
)

type Service interface {
	GetTenderList(ctx context.Context, serviceType string) (*[]Tender, error)
}

type tenderServices struct {
	log *slog.Logger
	db  *sql.DB
}

func NewService(log *slog.Logger, db *sql.DB) *tenderServices {
	return &tenderServices{
		log: log,
		db:  db,
	}
}

func (s *tenderServices) GetTenderList(serviceType string) ([]*Tender, error) {
	const op = "service.tender.GetTenderList"
	var tenders []*Tender
	var q string

	if len(serviceType) > 0 {
		st := strings.Split(serviceType, ",")
		q = fmt.Sprintf("SELECT %s FROM tender WHERE tender.service_type IN ('%s')", strings.Join(Columns, ","), strings.Join(st, "', '"))
	} else {
		q = fmt.Sprintf("SELECT %s FROM tender", strings.Join(Columns, ","))
	}

	rows, err := s.db.Query(q)

	switch {
	case err == sql.ErrNoRows:
		return []*Tender{}, fmt.Errorf("%s:%w", op, err)
	case err != nil:
		return []*Tender{}, fmt.Errorf("%s:%w", op, err)
	}

	defer rows.Close()

	//var responseIds []string
	for rows.Next() {
		tender, err := s.tenderFromRows(rows)
		if err != nil {
			return tenders, fmt.Errorf("%s:%w", op, err)
		}

		//tender.Employee = employee
		//if tender.Employee == nil {
		//	responseIds = append(responseIds, tender.responseId)
		//}
		tenders = append(tenders, tender)
	}

	if err := rows.Err(); err != nil {
		return tenders, fmt.Errorf("%s:%w", op, err)
	}

	// TODO: get Employee and Response Organization
	//if len(responseIds) > 0 {
	//	for i, tender := range tenders {
	//		userSet, err := getResponseTender(tender[i]["responseId"])
	//		if err != nil {
	//			return tenders, fmt.Errorf("%s:%w", op, err)
	//		}
	//		tenders[i].User = userSet[tender.responseId]
	//	}
	//}

	return tenders, nil
}

func (s *tenderServices) tenderFromRows(row *sql.Rows) (*Tender, error) {
	var t Tender
	err := row.Scan(&t.Id, &t.Name, &t.Description, &t.Status, &t.ResponseId, &t.ResponseId)
	return &t, err
}

func (s *tenderServices) getResponseTender(id int) ([]employee.Employee, error) {
	const op = "service.tender.getResponseTender"
	if id <= 0 {
		return []employee.Employee{}, fmt.Errorf("%s:%w", op, "Parametr id is empty")
	}

	return []employee.Employee{}, nil
}
