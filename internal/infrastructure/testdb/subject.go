package testdb

import (
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"math/rand"
	"time"
)

// 科目数据
func SeedSubject(db *pg.DB, users []*model.User) ([]*model.Subject, error) {
	subjects := []*model.Subject{
		{Name: "科目一", UpdatedAt: time.Now(), CreatedAt: time.Now()},
		{Name: "科目二", UpdatedAt: time.Now(), CreatedAt: time.Now()},
		{Name: "科目三", UpdatedAt: time.Now(), CreatedAt: time.Now()},
	}
	for _, subject := range subjects {
		subject.CreatedById = users[rand.Intn(len(users))].Id
		_, err := db.Model(subject).Returning("*").Insert()
		if err != nil {
			return nil, err
		}
	}
	return subjects, nil
}
