package testdb

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/database"
)

func New() (*pg.DB, error) {
	db, err := database.New(context.Background(), &pg.Options{
		Addr:     "localhost:5432",
		User:     "postgres",
		Password: "1234",
		Database: "example2",
	})
	if err != nil {
		return nil, err
	}

	err = Truncate(db)
	if err != nil {
		return nil, err
	}
	return db, err
}

func Truncate(db *pg.DB) error {
	stmt := `TRUNCATE TABLE "user", class, subject, learning_material`
	_, err := db.Exec(stmt)
	return err
}

func Drop(db *pg.DB) error {
	stmt := `DROP TABLE "user", class, subject, learning_material`
	_, err := db.Exec(stmt)
	return err
}
