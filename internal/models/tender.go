package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type Tender struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"service_type"`
	Status      string `json:"status"`
	ResponseId  string `json:"responsible"`

	Employee *Employee
}

var tenderColumns = []string{"id", "name", "description", "service_type", "status", "responsible"}

func tenderFromRows(row *sql.Rows) (*Tender, error) {
	var t Tender
	err := row.Scan(&t.Id, &t.Name, &t.Description, &t.ServiceType, &t.Status, &t.ResponseId)
	return &t, err
}

func GetTenderList(db *sql.DB, serviceType string) ([]*Tender, error) {
	const op = "model.tender.GetTenderList"
	var tenders []*Tender
	var q string

	if len(serviceType) > 0 {
		st := strings.Split(serviceType, ",")
		q = fmt.Sprintf("SELECT %s FROM tender WHERE tender.service_type IN ('%s')", strings.Join(tenderColumns, ","), strings.Join(st, "', '"))
	} else {
		q = fmt.Sprintf("SELECT %s FROM tender", strings.Join(tenderColumns, ","))
	}

	rows, err := db.Query(q)

	switch {
	case err == sql.ErrNoRows:
		return []*Tender{}, fmt.Errorf("%s:%w", op, err)
	case err != nil:
		return []*Tender{}, fmt.Errorf("%s:%w", op, err)
	}

	defer rows.Close()

	//var responseIds []string
	for rows.Next() {
		tender, err := tenderFromRows(rows)
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

func getResponseTender(id int) ([]Employee, error) {
	const op = "model.tender.getResponseTender"
	if id <= 0 {
		return []Employee{}, fmt.Errorf("%s:%w", op, "Parametr id is empty")
	}

	return []Employee{}, nil
}
