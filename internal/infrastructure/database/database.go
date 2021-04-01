package database

import (
	"context"
	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/setting"
	"github.com/xuxusheng/time-frequency-be/internal/model"
)

func New(setting *setting.Setting) (*pg.DB, error) {
	ctx := context.Background()

	db := pg.Connect(&pg.Options{
		Addr:     setting.DB.Host,
		User:     setting.DB.User,
		Password: setting.DB.Password,
		Database: setting.DB.Database,
	})

	if err := db.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "数据库连接失败")
	}

	if setting.Server.Mode == "debug" {
		db.AddQueryHook(pgdebug.DebugHook{
			Verbose: true,
		})
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
