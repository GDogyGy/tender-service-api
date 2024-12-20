package organizationResponsible

import (
	"TenderServiceApi/internal/model"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type OrganizationResponsible struct {
	db *sqlx.DB
}

func NewOrganizationResponsibleFacade(db *sqlx.DB) *OrganizationResponsible {
	return &OrganizationResponsible{db: db}
}

type argument struct {
	CreatorUsername string
	OrganizationID  string
}

// Зачем вообще фасад понадобился: приходит в пост не полная модель тендера с именем пользователя и id организации.
// Нужно перед созданием проверить имеет ли права юзер от этой организации. И желательно вытащить id usera по username и по двум параметрам кинуть запрос
// в organization_responsible чтобы в tender вставить id
// TODO: Правильно ли положил логику?
// TODO не нравится решение args можно все передавать, нужно обсудить, как решается изящнеe, когда нужно передавать несколько аргументов
func (o *OrganizationResponsible) Fetch(ctx context.Context, args []byte) (model.OrganizationResponsible, error) {
	const op = "repository.facade.organizationResponsible.fetch"
	var a argument
	var result model.OrganizationResponsible

	err := json.Unmarshal(args, &a)

	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	q := fmt.Sprintf(`SELECT organization_responsible.id, organization_responsible.organization_id, organization_responsible.user_id FROM organization_responsible left join employee on employee.id = organization_responsible.user_id left join organization on organization.id = organization_responsible.organization_id WHERE employee.username = $1 AND organization.id = $2`)

	row := o.db.QueryRowxContext(ctx, q, a.CreatorUsername, a.OrganizationID)

	err = row.Err()

	if errors.Is(err, sql.ErrNoRows) {
		return result, model.NotFound
	}

	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	result, err = o.fromRow(row)

	if errors.Is(err, sql.ErrNoRows) {
		return result, fmt.Errorf("%s:%w", op, err)
	}

	if err != nil {
		return result, err
	}
	
	// TODO в тз нужно уточнить об organizationId  {
	//
	//    "name": "Тендер 1",
	//
	//    "description": "Описание тендера",
	//
	//    "serviceType": "Construction",
	//
	//    "status": "Open",
	//
	//    "organizationId": 1,
	//
	//    "creatorUsername": "user1"
	//
	//  } organizationId тут UUID же?

	return result, nil
}

func (o *OrganizationResponsible) fromRow(rows *sqlx.Row) (model.OrganizationResponsible, error) {
	var r row
	err := rows.StructScan(&r)
	m := r.toModel()
	return m, err
}
