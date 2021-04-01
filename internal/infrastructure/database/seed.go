package database

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/utils"
	"time"
)

// 如果数据库中没有用户的话，初始化一个 admin 用户
func seedAdmin(ctx context.Context, db *pg.DB) error {
	is, err := db.ModelContext(ctx, &model.User{}).Exists()
	if err != nil {
		return err
	}
	if !is {
		hash, err := utils.EncodePwd("admin")
		if err != nil {
			return err
		}

		_, err = db.ModelContext(ctx, &model.User{
			Name:      "admin",
			NickName:  "管理员",
			Phone:     "12345678901",
			Email:     "admin@admin.com",
			Role:      "teacher",
			IsAdmin:   true,
			Password:  hash,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}).Insert()
		if err != nil {
			return err
		}
	}
	return nil
}
