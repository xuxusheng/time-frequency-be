package global

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-playground/validator/v10"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/setting"
)

var (
	// 配置配置项
	Setting *setting.Setting

	// 全局通用的数据库实例
	DB *pg.DB

	Validator = validator.New()
)
