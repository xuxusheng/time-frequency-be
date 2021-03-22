package dao

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/assert"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/testdb"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"math/rand"
	"testing"
	"time"
)

func prepareLearningMaterial(t *testing.T, db *pg.DB) ([]*model.LearningMaterial, []*model.Subject, []*model.User) {
	users, err := testdb.SeedUser(db)
	if err != nil || len(users) == 0 {
		t.Fatalf("准备用户数据失败：%v", err)
	}

	subjects, err := testdb.SeedSubject(db, users)
	if err != nil || len(subjects) == 0 {
		t.Fatalf("准备科目数据失败：%v", err)
	}

	learningMaterials, err := testdb.SeedLearningMaterial(db, users, subjects)
	if err != nil || len(learningMaterials) == 0 {
		t.Fatalf("准备学习资料数据失败：%v", err)
	}
	return learningMaterials, subjects, users
}

func TestLearningMaterialDao_Create(t *testing.T) {
	pls, pSubjects, pUsers := prepareLearningMaterial(t, db)
	dao := NewLearningMaterial(db)

	createdById := pUsers[rand.Intn(len(pUsers))].Id
	subjectId := pSubjects[rand.Intn(len(pSubjects))].Id

	t.Run("资料名称重复", func(t *testing.T) {
		for _, pl := range pls {
			s := time.Now().String()
			description := s + "description"
			md5 := s + "md5"
			filePath := s + "filePath"
			l, err := dao.Create(context.Background(), createdById, subjectId, pl.Name, description, md5, filePath)
			assert.EqualError(t, err, "ERROR #23505 duplicate key value violates unique constraint \"learning_material_name_key\"")
			assert.Nil(t, l)
		}
	})

	t.Run("资料名称为空", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			l, err := dao.Create(context.Background(), createdById, subjectId, "", s, s, s)
			assert.NotNil(t, err)
			assert.Nil(t, l)
		}
	})

	t.Run("资料描述为空", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			l, err := dao.Create(context.Background(), createdById, subjectId, s, "", s, s)
			assert.Nil(t, err)
			assert.NotZero(t, l)
		}
	})

	t.Run("正常创建", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			name := s + "name"
			description := s + "description"
			md5 := s + "md5"
			filePath := s + "filePath"
			l, err := dao.Create(context.Background(), createdById, subjectId, name, description, md5, filePath)
			if assert.Nil(t, err) {
				assert.NotZero(t, l.Id)
				assert.Equal(t, name, l.Name)
				assert.Equal(t, description, l.Description)
				assert.Equal(t, createdById, l.CreatedById)
				assert.Equal(t, subjectId, l.SubjectId)
				assert.Equal(t, md5, l.Md5)
				assert.Equal(t, filePath, l.FilePath)
			}
		}
	})

	_ = testdb.Truncate(db)
}

func TestLearningMaterialDao_Get(t *testing.T) {
	pls, _, _ := prepareLearningMaterial(t, db)
	dao := NewLearningMaterial(db)

	t.Run("正常获取", func(t *testing.T) {
		for _, pl := range pls {
			l, err := dao.Get(context.Background(), pl.Id)
			if assert.Nil(t, err) {
				assert.Equal(t, pl, l)
			}
		}

	})

	t.Run("获取不存在的资料", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			l, err := dao.Get(context.Background(), rand.Intn(100)*10000)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, l)
		}
	})

	_ = testdb.Truncate(db)
}

func TestLearningMaterialDao_Update(t *testing.T) {
	pls, _, pUsers := prepareLearningMaterial(t, db)
	dao := NewLearningMaterial(db)

	updatedById := pUsers[rand.Intn(len(pUsers))].Id

	t.Run("资料名称重复", func(t *testing.T) {
		for i := 0; i < len(pls)-1; i++ {
			current := pls[i]
			next := pls[i+1]
			l, err := dao.Update(context.Background(), current.Id, updatedById, next.Name, time.Now().String())
			assert.NotNil(t, err)
			assert.Nil(t, l)
		}
	})

	t.Run("资料名称为空", func(t *testing.T) {
		for _, pl := range pls {
			l, err := dao.Update(context.Background(), pl.Id, updatedById, "", time.Now().String())
			assert.NotNil(t, err)
			assert.Nil(t, l)
		}
	})

	t.Run("资料描述为空", func(t *testing.T) {
		for _, pl := range pls {
			l, err := dao.Update(context.Background(), pl.Id, updatedById, time.Now().String(), "")
			assert.Nil(t, err)
			assert.Zero(t, l.Description)
		}
	})

	t.Run("修改人ID为零值", func(t *testing.T) {
		for _, pl := range pls {
			s := time.Now().String()
			l, err := dao.Update(context.Background(), pl.Id, 0, s, s)
			assert.NotNil(t, err)
			assert.Nil(t, l)
		}
	})

	t.Run("正常修改", func(t *testing.T) {
		for _, pl := range pls {
			s := time.Now().String()
			name := s + "name"
			description := s + "description"
			l, err := dao.Update(context.Background(), pl.Id, updatedById, name, description)
			if assert.Nil(t, err) {
				assert.Equal(t, name, l.Name)
				assert.Equal(t, description, l.Description)
				assert.Greater(t, l.UpdatedAt.UnixNano(), pl.UpdatedAt.UnixNano())
			}
		}
	})

	_ = testdb.Truncate(db)
}

func TestLearningMaterialDao_Delete(t *testing.T) {
	pLs, _, _ := prepareLearningMaterial(t, db)
	dao := NewLearningMaterial(db)

	t.Run("正常删除", func(t *testing.T) {
		at := assert.New(t)
		for _, pL := range pLs {
			err := dao.Delete(context.Background(), pL.Id)
			if at.Nil(err) {
				var l model.LearningMaterial
				err = db.Model(&l).Where("id = ?", pL.Id).Select()
				at.Zero(l)
				at.Equal(pg.ErrNoRows, err)
			}
		}
	})

	t.Run("删除不存在的资料", func(t *testing.T) {
		for i := 0; i < 0; i++ {
			err := dao.Delete(context.Background(), rand.Intn(100)*10000)
			assert.Nil(t, err)
		}
	})

	_ = testdb.Truncate(db)
}
