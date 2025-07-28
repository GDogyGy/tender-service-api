package client

import (
	"context"
	"os/signal"
	"path"
	"runtime"
	"syscall"

	"TenderServiceApi/internal/storage/postgres"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jmoiron/sqlx"
)

func InitStorageDB() (*sqlx.DB, error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := postgres.New(ctx, "postgres://root:123@localhost:5434/TenderApiTest?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db.Db, nil
}

func DeleteAllCollections(ctx context.Context, db *sqlx.DB) error {
	r := db.QueryRowxContext(ctx, `TRUNCATE organization_responsible, tender, bids, employee, organization, bid_feedback, decision_bid RESTART IDENTITY`)
	err := r.Err()
	if err != nil {
		return err
	}
	return nil
}

func PrepareDB(ctx context.Context, db *sqlx.DB) error {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := path.Dir(filename)
	fixtureDir := path.Join(currentDir, "..", "fixtures")

	fixtures, err := testfixtures.New(
		testfixtures.Database(db.DB), // You database connection
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			fixtureDir+"/organization.yml",
			fixtureDir+"/employee.yml",
			fixtureDir+"/organization_responsible.yml",
			fixtureDir+"/tender.yml",
			fixtureDir+"/bids.yml",
			fixtureDir+"/bid_feedback.yml",
			fixtureDir+"/decision_bid.yml",
		),
	)
	if err != nil {
		return err
	}

	return fixtures.Load()
}
