package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/app"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/database"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/setting"
	"log"
)

// @title 时频培训系统 API 接口文档
// @version 1.0
// @contact.name 许盛
// @contact.email 20691718@qq.com
// @host localhost:8080
func main() {
	var err error

	global.Setting, err = setupSetting()
	if err != nil {
		log.Fatalf("初始化全局配置项失败：%v", err)
	}

	global.DB, err = setupDB(global.Setting)
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
func setupDB(s *setting.Setting) (*pg.DB, error) {
	db, err := database.New(s)
	if err != nil {
		return nil, err
	}
	return db, err
}
