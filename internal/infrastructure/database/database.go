package database

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
	"github.com/xuxusheng/time-frequency-be/internal/model"
)

func New(ctx context.Context, options *pg.Options) (*pg.DB, error) {
	db := pg.Connect(options)

	if err := db.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "连接失败")
	}

	schemas := []interface{}{
		(*model.User)(nil),
		(*model.Class)(nil),
		(*model.Subject)(nil),
		(*model.LearningMaterial)(nil),
	}

	for _, schema := range schemas {
		err := db.ModelContext(ctx, schema).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return nil, errors.Wrap(err, "创建数据表失败")
		}
	}
	return db, nil
}
