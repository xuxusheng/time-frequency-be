package testdb

import (
	"github.com/go-pg/pg/v10"
)

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
