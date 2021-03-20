package service

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/assert"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/testdb"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
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
	pLms, pSubjects, pUsers := prepareLearningMaterial(t, db)
	svc := NewLearningMaterial(lmDao)

	t.Run("资料名称重复", func(t *testing.T) {
		for _, pLm := range pLms {
			createdById := pUsers[rand.Intn(len(pUsers))].Id
			subjectId := pSubjects[rand.Intn(len(pSubjects))].Id
			s := time.Now().String()
			_, err := svc.Create(context.Background(), createdById, subjectId, pLm.Name, s, s, s)
			assert.Equal(t, cerror.BadRequest.WithMsg("资料名称已存在"), err)
		}
	})

	t.Run("字段为空的情况", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			name := s + "name"
			desc := s + "desc"
			createdById := pUsers[rand.Intn(len(pUsers))].Id
			subjectId := pSubjects[rand.Intn(len(pSubjects))].Id
			md5 := s + "md5"
			filePath := s + "filePath"
			t.Run("名称为空", func(t *testing.T) {
				lm, err := svc.Create(context.Background(), createdById, subjectId, "", desc, md5, filePath)
				assert.NotNil(t, err)
				assert.Nil(t, lm)
			})
			t.Run("描述为空", func(t *testing.T) {
				lm, err := svc.Create(context.Background(), createdById, subjectId, name, "", md5, filePath)
				if assert.Nil(t, err) {
					assert.Equal(t, name, lm.Name)
					assert.Zero(t, lm.Description)
					assert.Equal(t, createdById, lm.CreatedById)
					assert.Equal(t, subjectId, lm.SubjectId)
					assert.Equal(t, md5, lm.Md5)
					assert.Equal(t, filePath, lm.FilePath)
				}
			})
			t.Run("创建人为空", func(t *testing.T) {
				lm, err := svc.Create(context.Background(), 0, subjectId, name, desc, md5, "")
				assert.NotNil(t, err)
				assert.Nil(t, lm)
			})

			t.Run("所属科目为空", func(t *testing.T) {
				lm, err := svc.Create(context.Background(), createdById, 0, name, desc, md5, "")
				assert.NotNil(t, err)
				assert.Nil(t, lm)
			})
			t.Run("md5为空", func(t *testing.T) {
				lm, err := svc.Create(context.Background(), createdById, subjectId, name, desc, "", filePath)
				assert.NotNil(t, err)
				assert.Nil(t, lm)
			})
			t.Run("filePath为空", func(t *testing.T) {
				lm, err := svc.Create(context.Background(), createdById, subjectId, name, desc, md5, "")
				assert.NotNil(t, err)
				assert.Nil(t, lm)
			})

		}
	})

	t.Run("正常创建", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			name := s + "name"
			desc := s + "desc"
			createdById := pUsers[rand.Intn(len(pUsers))].Id
			subjectId := pSubjects[rand.Intn(len(pSubjects))].Id
			md5 := s + "md5"
			filePath := s + "filePath"
			lm, err := svc.Create(context.Background(), createdById, subjectId, name, desc, md5, filePath)
			if assert.Nil(t, err) {
				assert.Equal(t, name, lm.Name)
				assert.Equal(t, desc, lm.Description)
				assert.Equal(t, createdById, lm.CreatedById)
				assert.Equal(t, subjectId, lm.SubjectId)
				assert.Equal(t, md5, lm.Md5)
				assert.Equal(t, filePath, lm.FilePath)
			}

		}
	})

	truncateTable(db)
}

func TestLearningMaterialSvc_Get(t *testing.T) {
	pLms, _, _ := prepareLearningMaterial(t, db)
	svc := NewLearningMaterial(lmDao)

	t.Run("正常获取", func(t *testing.T) {
		for _, pLm := range pLms {
			lm, err := svc.Get(context.Background(), pLm.Id)
			if assert.Nil(t, err) {
				assert.Equal(t, pLm, lm)
			}
		}
	})

	t.Run("获取不存在的资料", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			lm, err := svc.Get(context.Background(), rand.Intn(100)*10000)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, lm)
		}
	})

	truncateTable(db)
}

func TestLearningMaterialSvc_Update(t *testing.T) {
	pLms, _, pUsers := prepareLearningMaterial(t, db)
	svc := NewLearningMaterial(lmDao)

	t.Run("班级名称重复", func(t *testing.T) {
		for i := 0; i < len(pLms)-1; i++ {
			current := pLms[i]
			next := pLms[i+1]
			updatedById := pUsers[rand.Intn(len(pUsers))].Id
			lm, err := svc.Update(context.Background(), current.Id, updatedById, next.Name, time.Now().String())
			assert.EqualError(t, err, "资料名称已存在")
			assert.Nil(t, lm)
		}
	})

	t.Run("资料名称为空", func(t *testing.T) {
		for _, pLm := range pLms {
			updatedById := pUsers[rand.Intn(len(pUsers))].Id
			lm, err := svc.Update(context.Background(), pLm.Id, updatedById, "", time.Now().String())
			assert.NotNil(t, err)
			assert.Nil(t, lm)
		}
	})

	t.Run("修改一个不存在的资料", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			updatedById := pUsers[rand.Intn(len(pUsers))].Id
			lm, err := svc.Update(context.Background(), rand.Intn(100)*10000, updatedById, s, s)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, lm)
		}
	})

	t.Run("正常修改", func(t *testing.T) {
		for _, pLm := range pLms {
			s := time.Now().String()
			name := s + "name"
			desc := s + "desc"
			updatedById := pUsers[rand.Intn(len(pUsers))].Id
			lm, err := svc.Update(context.Background(), pLm.Id, updatedById, name, desc)
			if assert.Nil(t, err) {
				assert.Equal(t, name, lm.Name)
				assert.Equal(t, desc, lm.Description)
				assert.Greater(t, lm.UpdatedAt.UnixNano(), pLm.UpdatedAt.UnixNano())
			}
		}
	})

	truncateTable(db)
}

func TestLearningMaterialSvc_Delete(t *testing.T) {
	pLms, _, _ := prepareLearningMaterial(t, db)
	svc := NewLearningMaterial(lmDao)

	t.Run("正常删除", func(t *testing.T) {
		for _, pLm := range pLms {
			err := svc.Delete(context.Background(), pLm.Id)
			if assert.Nil(t, err) {
				var lm model.LearningMaterial
				err = db.Model(&lm).Where("id = ?", pLm.Id).Select()
				assert.Zero(t, lm)
				assert.Equal(t, pg.ErrNoRows, err)
			}
		}
	})

	t.Run("删除不存在的资料", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			err := svc.Delete(context.Background(), rand.Intn(100)*10000)
			assert.Nil(t, err)
		}
	})

	truncateTable(db)
}

func TestLearningMaterialSvc_IsNameExist(t *testing.T) {
	pLms, _, _ := prepareLearningMaterial(t, db)
	svc := NewLearningMaterial(lmDao)

	t.Run("排除当前资料后，查找当前资料的名称", func(t *testing.T) {
		for _, pLm := range pLms {
			is, err := svc.IsNameExist(context.Background(), pLm.Name, pLm.Id)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	t.Run("排除当前资料后，查找其他资料名称", func(t *testing.T) {
		for i := 0; i < len(pLms)-1; i++ {
			current := pLms[i]
			next := pLms[i+1]
			is, err := svc.IsNameExist(context.Background(), next.Name, current.Id)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找存在的资料名称", func(t *testing.T) {
		for _, pLm := range pLms {
			is, err := svc.IsNameExist(context.Background(), pLm.Name, 0)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找不存在的资料名称", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			is, err := svc.IsNameExist(context.Background(), time.Now().String(), 0)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	truncateTable(db)
}
