package dao

import (
	"github.com/go-pg/pg/v10"
)

type Dao struct {
	engine *pg.DB
}
