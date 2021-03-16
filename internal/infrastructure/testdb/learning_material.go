package testdb

import (
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"math/rand"
	"time"
)

func SeedLearningMaterial(db *pg.DB, users []*model.User, subjects []*model.Subject) ([]*model.LearningMaterial, error) {

	lms := []*model.LearningMaterial{
		{Name: "资料一", Description: "资料一", Md5: "md5", FilePath: "filePath"},
		{Name: "资料二", Description: "资料二", Md5: "md5", FilePath: "filePath"},
		{Name: "资料三", Description: "资料三", Md5: "md5", FilePath: "filePath"},
		{Name: "资料四", Description: "资料四", Md5: "md5", FilePath: "filePath"},
	}
	for _, l := range lms {
		now := time.Now()
		l.UpdatedAt = now
		l.CreatedAt = now
		l.CreatedById = users[rand.Intn(len(users))].Id
		l.UpdatedById = users[rand.Intn(len(users))].Id
		l.SubjectId = subjects[rand.Intn(len(subjects))].Id
		_, err := db.Model(l).Returning("*").Insert()
		if err != nil {
			return nil, err
		}
	}
	return lms, nil
}
