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

func prepareSubject(t *testing.T, db *pg.DB) ([]*model.Subject, []*model.User) {
	users, err := testdb.SeedUser(db)
	if err != nil || len(users) == 0 {
		t.Fatalf("准备用户数据失败：%v", err)
	}
	subjects, err := testdb.SeedSubject(db, users)
	if err != nil || len(subjects) == 0 {
		t.Fatalf("准备科目数据失败：%v", err)
	}
	return subjects, users
}

func TestSubjectDao_Create(t *testing.T) {
	pSubjects, pUsers := prepareSubject(t, db)
	createdById := pUsers[0].Id
	dao := NewSubject(db)

	t.Run("名称重复", func(t *testing.T) {
		for _, pSubject := range pSubjects {
			subject, err := dao.Create(context.Background(), createdById, pSubject.Name, time.Now().String())
			assert.NotNil(t, err)
			assert.Nil(t, subject)
		}
	})

	t.Run("科目名称为空", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			subject, err := dao.Create(context.Background(), createdById, "", time.Now().String())
			assert.NotNil(t, err)
			assert.Nil(t, subject)
		}
	})

	t.Run("正常创建", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			name := s + "name"
			description := s + "description"
			subject, err := dao.Create(context.Background(), createdById, name, description)
			if assert.Nil(t, err) {
				assert.NotZero(t, subject.Id)
				assert.Equal(t, name, subject.Name)
				assert.Equal(t, description, subject.Description)
				assert.NotNil(t, subject.CreatedAt)
				assert.NotNil(t, subject.UpdatedAt)
			}
		}
	})

	truncateTable(db)
}

func TestSubjectDao_Get(t *testing.T) {
	pSubjects, _ := prepareSubject(t, db)
	dao := NewSubject(db)

	t.Run("正常获取", func(t *testing.T) {
		for _, pSubject := range pSubjects {
			subject, err := dao.Get(context.Background(), pSubject.Id)
			if assert.Nil(t, err) {
				assert.Equal(t, pSubject, subject)
			}
		}
	})

	t.Run("获取不存在的科目", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			subject, err := dao.Get(context.Background(), rand.Intn(100)*10000)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, subject)
		}
	})
	truncateTable(db)
}

func TestSubjectDao_Update(t *testing.T) {
	pSubjects, _ := prepareSubject(t, db)
	dao := NewSubject(db)

	t.Run("科目名称重复", func(t *testing.T) {
		for i := 0; i < len(pSubjects)-1; i++ {
			current := pSubjects[i]
			next := pSubjects[i+1]
			subject, err := dao.Update(context.Background(), current.Id, next.Name, time.Now().String())
			assert.Nil(t, subject)
			assert.NotNil(t, err)
		}
	})

	t.Run("科目名称为空", func(t *testing.T) {
		for _, pSubject := range pSubjects {
			subject, err := dao.Update(context.Background(), pSubject.Id, "", time.Now().String())
			assert.NotNil(t, err)
			assert.Nil(t, subject)
		}
	})

	t.Run("修改一个不存在的科目", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			subject, err := dao.Update(context.Background(), rand.Intn(100)*10000, s, s)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, subject)
		}
	})

	t.Run("正常修改", func(t *testing.T) {
		for _, pSubject := range pSubjects {
			s := time.Now().String()
			name := s + "name"
			description := s + "description"
			subject, err := dao.Update(context.Background(), pSubject.Id, name, description)
			if assert.Nil(t, err) {
				assert.Equal(t, name, subject.Name)
				assert.Equal(t, description, subject.Description)
				assert.Greater(t, subject.UpdatedAt.UnixNano(), pSubject.UpdatedAt.UnixNano())
			}
		}
	})

	truncateTable(db)
}

func TestSubjectDao_Delete(t *testing.T) {
	pSubjects, _ := prepareSubject(t, db)
	dao := NewSubject(db)

	t.Run("正常删除", func(t *testing.T) {
		for _, pSubject := range pSubjects {
			err := dao.Delete(context.Background(), pSubject.Id)
			if assert.Nil(t, err) {
				var subject model.Subject
				err = db.Model(&subject).Where("id = ?", pSubject.Id).Select()
				assert.Zero(t, subject)
				assert.Equal(t, pg.ErrNoRows, err)
			}
		}
	})

	t.Run("删除不存在的科目", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			err := dao.Delete(context.Background(), rand.Intn(100)*10000)
			assert.Nil(t, err)
		}
	})

	truncateTable(db)
}
