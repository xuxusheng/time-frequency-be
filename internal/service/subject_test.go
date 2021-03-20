package service

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

func prepareSubject(t *testing.T, db *pg.DB) ([]*model.Subject, []*model.User) {
	users, err := testdb.SeedUser(db)
	if err != nil || len(users) == 0 {
		t.Fatalf("准备用户数据失败：%v", err)
	}
	subjects, err := testdb.SeedSubject(db, users)
	if err != nil || len(users) == 0 {
		t.Fatalf("准备科目数据失败：%v", err)
	}
	return subjects, users
}

func TestSubjectSvc_Create(t *testing.T) {
	pSubjects, pUsers := prepareSubject(t, db)
	svc := NewSubject(subjectDao)

	t.Run("科目名称重复", func(t *testing.T) {
		for _, pSubject := range pSubjects {
			createdById := pUsers[rand.Intn(len(pUsers))].Id
			_, err := svc.Create(context.Background(), createdById, pSubject.Name, time.Now().String())
			assert.EqualError(t, err, "科目名称已存在")
		}
	})

	t.Run("科目名称、描述为空", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			name := s + "name"
			description := s + "description"
			createdById := pUsers[rand.Intn(len(pUsers))].Id
			t.Run("名称为空", func(t *testing.T) {
				_, err := svc.Create(context.Background(), createdById, "", description)
				assert.NotNil(t, err)
			})
			t.Run("描述为空", func(t *testing.T) {
				subject, err := svc.Create(context.Background(), createdById, name, "")
				if assert.Nil(t, err) {
					assert.Equal(t, name, subject.Name)
					assert.Equal(t, createdById, subject.CreatedById)
					assert.Zero(t, subject.Description)
				}
			})
		}
	})

	t.Run("正常创建", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			name := s + "name"
			description := s + "description"
			createdById := pUsers[rand.Intn(len(pUsers))].Id
			subject, err := svc.Create(context.Background(), createdById, name, description)
			if assert.Nil(t, err) {
				assert.Equal(t, name, subject.Name)
				assert.Equal(t, description, subject.Description)
				assert.Equal(t, createdById, subject.CreatedById)
			}
		}
	})

	truncateTable(db)
}

func TestSubjectSvc_Get(t *testing.T) {
	pSubjects, _ := prepareSubject(t, db)
	svc := NewSubject(subjectDao)

	t.Run("正常获取", func(t *testing.T) {
		for _, pSubject := range pSubjects {
			subject, err := svc.Get(context.Background(), pSubject.Id)
			if assert.Nil(t, err) {
				assert.Equal(t, pSubject, subject)
			}
		}
	})

	t.Run("获取不存在的科目", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			subject, err := svc.Get(context.Background(), rand.Intn(100)*10000)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, subject)
		}
	})

	truncateTable(db)
}

func TestSubjectSvc_Update(t *testing.T) {
	pSubjects, _ := prepareSubject(t, db)
	svc := NewSubject(subjectDao)

	t.Run("科目名称重复", func(t *testing.T) {
		for i := 0; i < len(pSubjects)-1; i++ {
			current := pSubjects[i]
			next := pSubjects[i+1]
			subject, err := svc.Update(context.Background(), current.Id, next.Name, time.Now().String())
			assert.EqualError(t, err, "科目名称已存在")
			assert.Nil(t, subject)
		}
	})

	t.Run("科目名称为空", func(t *testing.T) {
		for _, pSubject := range pSubjects {
			subject, err := svc.Update(context.Background(), pSubject.Id, "", time.Now().String())
			assert.NotNil(t, err)
			assert.Nil(t, subject)
		}
	})

	t.Run("修改一个不存在的科目", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			subject, err := svc.Update(context.Background(), rand.Intn(100)*10000, s, s)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, subject)
		}
	})

	t.Run("正常修改", func(t *testing.T) {
		for _, pSubject := range pSubjects {
			s := time.Now().String()
			name := s + "name"
			description := s + "description"
			subject, err := svc.Update(context.Background(), pSubject.Id, name, description)
			if assert.Nil(t, err) {
				assert.Equal(t, name, subject.Name)
				assert.Equal(t, description, subject.Description)
				assert.Greater(t, subject.UpdatedAt.UnixNano(), pSubject.UpdatedAt.UnixNano())
			}

		}
	})

	truncateTable(db)
}

func TestSubjectSvc_Delete(t *testing.T) {
	pSubjects, _ := prepareSubject(t, db)
	svc := NewSubject(subjectDao)

	t.Run("正常删除", func(t *testing.T) {
		for _, pSubject := range pSubjects {
			err := svc.Delete(context.Background(), pSubject.Id)
			if assert.Nil(t, err) {
				var subject model.Subject
				err = db.Model(&subject).Where("id = ?", pSubject.Id).Select()
				assert.Zero(t, subject)
				assert.Equal(t, pg.ErrNoRows, err)
			}
		}
	})

	t.Run("删除一个不存在的科目", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			err := svc.Delete(context.Background(), rand.Intn(100)*10000)
			assert.Nil(t, err)
		}
	})

	truncateTable(db)
}

func TestSubjectSvc_IsNameExist(t *testing.T) {
	pSubjects, _ := prepareSubject(t, db)
	svc := NewSubject(subjectDao)

	t.Run("排除当前科目后，查找当前科目名称", func(t *testing.T) {
		for _, pSubject := range pSubjects {
			is, err := svc.IsNameExist(context.Background(), pSubject.Name, pSubject.Id)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	t.Run("排除当前科目后，查找其他科目名称", func(t *testing.T) {
		for i := 0; i < len(pSubjects)-1; i++ {
			current := pSubjects[i]
			next := pSubjects[i+1]
			is, err := svc.IsNameExist(context.Background(), next.Name, current.Id)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找存在的名称", func(t *testing.T) {
		for _, pSubject := range pSubjects {
			is, err := svc.IsNameExist(context.Background(), pSubject.Name, 0)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找不存在的名称", func(t *testing.T) {
		for i := 0; i < len(pSubjects)-1; i++ {
			is, err := svc.IsNameExist(context.Background(), time.Now().String(), 0)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	truncateTable(db)
}
