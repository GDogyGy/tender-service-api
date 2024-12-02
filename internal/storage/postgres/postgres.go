package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	Db *sqlx.DB
}

func New(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Connect("postgres", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}
	stmt, err := db.PrepareContext(ctx, `SELECT EXISTS (SELECT FROM public.tender)`)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{Db: db}, nil
}

func (s *Storage) Close() {
	_ = s.Db.Close()
}
