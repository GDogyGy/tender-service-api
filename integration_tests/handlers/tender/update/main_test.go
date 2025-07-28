//go:build integration
// +build integration

package test

import (
	"os"
	"testing"

	"TenderServiceApi/integration_tests/helpers/db/client"
	"github.com/jmoiron/sqlx"
)

var (
	Db *sqlx.DB
)

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		os.Exit(1)
	}

	exitCode := m.Run()

	if err := tearDown(); err != nil {
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func setup() error {
	var err error

	if Db, err = client.InitStorageDB(); err != nil {
		return err
	}

	return nil
}

func tearDown() error {
	if err := Db.Close(); err != nil {
		return err
	}
	return nil
}
