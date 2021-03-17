package service

import (
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/testdb"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"testing"
)

func prepareLearningMaterial(t *testing.T, db *pg.DB) ([]*model.LearningMaterial, []*model.Subject, []*model.User) {
	users, err := testdb.SeedUser(db)
	if err != nil || len(users) == 0 {
		t.Fatalf("准备用户数据失败：%v", err)
	}
	subjects, err := testdb.SeedSubject(db, users)
	if err != nil || len(users) == 0 {
		t.Fatalf("准备科目数据失败：%v", err)
	}
	lms, err := testdb.SeedLearningMaterial(db, users, subjects)
	if err != nil {
		t.Fatalf("准备资料数据失败：%v", err)
	}
	return lms, subjects, users
}

func TestLearningMaterialSvc_Create(t *testing.T) {
	//pLms, _, _ := prepareLearningMaterial(t, db)
	//svc := NewLearningMaterialSvc(lmDao)

	t.Run("", func(t *testing.T) {

	})

	truncateTable(db)
}

func TestLearningMaterialSvc_Get(t *testing.T) {
	//pLms, _, _ := prepareLearningMaterial(t, db)
	//svc := NewLearningMaterialSvc(lmDao)

	t.Run("", func(t *testing.T) {

	})

	truncateTable(db)
}

func TestLearningMaterialSvc_Update(t *testing.T) {
	//pLms, _, _ := prepareLearningMaterial(t, db)
	//svc := NewLearningMaterialSvc(lmDao)

	t.Run("", func(t *testing.T) {

	})

	truncateTable(db)
}

func TestLearningMaterialSvc_Delete(t *testing.T) {
	//pLms, _, _ := prepareLearningMaterial(t, db)
	//svc := NewLearningMaterialSvc(lmDao)

	t.Run("", func(t *testing.T) {

	})

	truncateTable(db)
}

func TestLearningMaterialSvc_IsNameExist(t *testing.T) {
	//pLms, _, _ := prepareLearningMaterial(t, db)
	//svc := NewLearningMaterialSvc(lmDao)

	t.Run("", func(t *testing.T) {

	})

	truncateTable(db)
}
