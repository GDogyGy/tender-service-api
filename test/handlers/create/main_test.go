//go:build integration
// +build integration

package test

import (
	"TenderServiceApi/internal/handlers/tender"
	helper "TenderServiceApi/test"
	"context"
	"github.com/jmoiron/sqlx"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite

	DB  *sqlx.DB
	Log *Mocklog
}

type Storage struct {
	Db *sqlx.DB
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	storage, err := sqlx.Connect("postgres", "postgres://root:123@localhost:5434/TenderApiTest?sslmode=disable")
	if err != nil {
		s.FailNow("Failed to connect to postgres", err)
	}

	// TODO: Тут импорт не получается моков организовать, чтобы handler завелся
	// TODO: Задача собрать в APITestSuite все что нужно для тестов (service, storage) и протестить хендлер
	logMock := tender.NewMocklog()
	s.DB = storage
	s.Log = slog
	err = helper.Migration("up", storage)
	if err != nil {
		s.FailNow("Failed to migrate postgres db", err)
	}

	_, err = storage.PrepareContext(context.TODO(), `SELECT EXISTS (SELECT FROM public.organization_responsible)`)
	if err != nil {
		s.Error(err)
	}

}

func (s *APITestSuite) TearDownSuite() {
	_ = s.DB.DB.Close()
}
