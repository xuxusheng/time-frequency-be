package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/pkg/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Model struct {
	ID        uint      `json:"id" gorm:"primarykey" example:"1"`                   // 唯一主键
	CreatedAt time.Time `json:"created_at" example:"2020-12-09T18:52:41.555+08:00"` // 记录创建时间
	UpdatedAt time.Time `json:"updated_at" example:"2020-12-09T18:52:41.555+08:00"` // 记录最后更新时间
}

// 新建数据库连接
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	)

	dbConfig := gorm.Config{}

	if global.ServerSetting.RunMode == gin.DebugMode {
		dbConfig.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      true,        // Disable color
			},
		)
	}

	// todo 后面加一个 DB_DSN 的配置项
	db, err := gorm.Open(
		mysql.Open(dsn),
		&dbConfig,
	)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil
}
