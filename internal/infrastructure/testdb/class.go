package testdb

import (
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"math/rand"
	"time"
)

// 班级数据
func SeedClass(db *pg.DB, users []*model.User) ([]*model.Class, error) {
	classes := []*model.Class{
		{Name: "一班", Description: "一班", UpdatedAt: time.Now(), CreatedAt: time.Now()},
		{Name: "二班", Description: "二班", UpdatedAt: time.Now(), CreatedAt: time.Now()},
		{Name: "三班", Description: "三班", UpdatedAt: time.Now(), CreatedAt: time.Now()},
		{Name: "四班", Description: "四班", UpdatedAt: time.Now(), CreatedAt: time.Now()},
		{Name: "五班", Description: "五班", UpdatedAt: time.Now(), CreatedAt: time.Now()},
	}
	for _, class := range classes {
		// 从 users 中随便取一个作为创建者
		class.CreatedById = users[rand.Intn(len(users))].Id
		_, err := db.Model(class).Returning("*").Insert()
		if err != nil {
			return nil, err
		}
	}
	return classes, nil
}
