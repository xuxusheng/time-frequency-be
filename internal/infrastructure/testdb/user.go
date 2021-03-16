package testdb

import (
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"time"
)

func SeedUser(db *pg.DB) ([]*model.User, error) {
	users := []*model.User{
		{Name: "张三", Phone: "1", Email: "1@qq.com", Password: "1234"},
		{Name: "李四", Phone: "2", Email: "2@qq.com", Password: "2345"},
		{Name: "王五", Phone: "3", Email: "3@qq.com", Password: "3456"},
		{Name: "胡六", Phone: "4", Email: "4@qq.com", Password: "4567"},
	}
	var lastUser *model.User
	for _, user := range users {
		if lastUser != nil {
			user.CreatedById = lastUser.Id
		}
		now := time.Now()
		user.CreatedAt = now
		user.UpdatedAt = now
		_, err := db.Model(user).Returning("*").Insert()
		if err != nil {
			return nil, err
		}
		// 把 createdBy 取出来，和 dao 中的方法保持一致
		err = db.Model(user).WherePK().Relation("CreatedBy").Select()
		if err != nil {
			return nil, err
		}
		lastUser = user
	}
	return users, nil
}
