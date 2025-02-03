package helper

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"path/filepath"
)

// TODO: Затем в main_test.go поднять внутри сервер чтобы апи тестовое работало
// TODO: протестить создание handler tender create
func MigrationRefresh(cmd string, db *sqlx.DB) error {
	if err := Migration("reset", db); err != nil {
		return err
	}

	if err := Migration("up", db); err != nil {
		return err
	}
	return nil
}

func Migration(cmd string, db *sqlx.DB) error {
	// TODO: Эта штука валидна? и если нет то как решить через filepath. OS.
	dirname := "../../../" + string(filepath.Separator) + "/internal/storage/migrations"

	if err := goose.Run(cmd, db.DB, dirname); err != nil {
		return err
	}
	return nil
}
