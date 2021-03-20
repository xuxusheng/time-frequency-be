package main

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/app"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/database"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/setting"
	"log"
)

func main() {
	var err error

	global.Setting, err = setupSetting()
	if err != nil {
		log.Fatalf("初始化全局配置项失败：%v", err)
	}

	global.DB, err = setupDB(global.Setting.DB)
	if err != nil {
		log.Fatalf("初始化数据库连接失败：%v", err)
	}

	a := app.New()

	a.Listen(":8080")
}

// 初始化全局配置项
func setupSetting() (*setting.Setting, error) {
	s, err := setting.New()
	if err != nil {
		return nil, err
	}

	err = s.Init()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// 初始化数据库连接
func setupDB(s *setting.DB) (*pg.DB, error) {
	db, err := database.New(context.Background(), &pg.Options{
		Database: s.Database,
		Addr:     s.Host,
		User:     s.User,
		Password: s.Password,
	})
	if err != nil {
		return nil, err
	}
	return db, err
}
