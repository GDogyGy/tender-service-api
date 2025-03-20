//go:build integration
// +build integration

package test

import (
	"fmt"
	"os"
	"testing"

	"TenderServiceApi/integration_tests/db/client"
	"github.com/jmoiron/sqlx"
)

var (
	Db *sqlx.DB
)

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		os.Exit(1)
	}

	fmt.Println("testing")

	exitCode := m.Run()

	if err := tearDown(); err != nil {
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func setup() error {
	var err error

	if Db, err = client.InitStorageDB(); err != nil {
		fmt.Println("Ошибка подключения", err)
		return err
	}
	fmt.Println(err)
	return nil
}

func tearDown() error {
	if err := Db.Close(); err != nil {
		return err
	}
	return nil
}
