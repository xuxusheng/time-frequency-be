package pg_dao

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/xuxusheng/time-frequency-be/internal/pg_model"
	"github.com/xuxusheng/time-frequency-be/pkg/setting"
	"log"
)

type Dao struct {
	engine *pg.DB
}

func NewPGEngine(setting *setting.PGSettingS) (*pg.DB, error) {

	db := pg.Connect(&pg.Options{
		Addr:     setting.Host,
		User:     setting.Username,
		Password: setting.Password,
		Database: setting.DBName,
	})

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}
	log.Println("pg 数据库连接成功")

	// 创建 table
	err := db.Model(&pg_model.User{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
		//Temp: true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
