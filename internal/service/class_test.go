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

func prepareClass(t *testing.T, db *pg.DB) ([]*model.Class, []*model.User) {
	users, err := testdb.SeedUser(db)
	if err != nil {
		t.Fatalf("准备用户数据失败：%v", err)
	}
	classes, err := testdb.SeedClass(db, users)
	if err != nil {
		t.Fatalf("准备班级数据失败：%v", err)
	}
	return classes, users
}

func TestClassSvc_Create(t *testing.T) {
	pClasses, pUsers := prepareClass(t, db)
	svc := NewClass(classDao)

	t.Run("班级名称重复", func(t *testing.T) {
		for _, pClass := range pClasses {
			createdById := pUsers[rand.Intn(len(pUsers))].Id
			class, err := svc.Create(context.Background(), createdById, pClass.Name, time.Now().String())
			assert.Equal(t, cerror.BadRequest.WithMsg("班级名称已存在"), err)
			assert.Nil(t, class)
		}
	})

	t.Run("班级名称、描述为空", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			name := s + "name"
			desc := s + "desc"
			createdById := pUsers[rand.Intn(len(pUsers))].Id
			t.Run("名称为空", func(t *testing.T) {
				class, err := svc.Create(context.Background(), createdById, "", desc)
				assert.NotNil(t, err)
				assert.Nil(t, class)
			})
			t.Run("描述为空", func(t *testing.T) {
				class, err := svc.Create(context.Background(), createdById, name, "")
				if assert.Nil(t, err) {
					assert.Equal(t, name, class.Name)
					assert.Zero(t, class.Description)
					assert.Equal(t, createdById, class.CreatedById)
				}
			})
		}
	})

	t.Run("正常创建", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			name := s + "name"
			desc := s + "desc"
			createdById := pUsers[rand.Intn(len(pUsers))].Id
			class, err := svc.Create(context.Background(), createdById, name, desc)
			if assert.Nil(t, err) {
				assert.Equal(t, name, class.Name)
				assert.Equal(t, desc, class.Description)
				assert.Equal(t, createdById, class.CreatedById)
			}
		}
	})

	truncateTable(db)
}

func TestClassSvc_Get(t *testing.T) {
	pClasses, _ := prepareClass(t, db)
	svc := NewClass(classDao)

	t.Run("正常获取", func(t *testing.T) {
		for _, pClass := range pClasses {
			class, err := svc.Get(context.Background(), pClass.Id)
			if assert.Nil(t, err) {
				assert.Equal(t, pClass, class)
			}
		}
	})

	t.Run("获取不存在的科目", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			class, err := svc.Get(context.Background(), rand.Intn(100)*10000)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, class)
		}
	})

	truncateTable(db)
}

func TestClassSvc_Update(t *testing.T) {
	pClasses, _ := prepareClass(t, db)
	svc := NewClass(classDao)

	t.Run("班级名称重复", func(t *testing.T) {
		for i := 0; i < len(pClasses)-1; i++ {
			current := pClasses[i]
			next := pClasses[i+1]
			class, err := svc.Update(context.Background(), current.Id, next.Name, time.Now().String())
			assert.EqualError(t, err, "班级名称已存在")
			assert.Nil(t, class)
		}
	})

	t.Run("班级名称为空", func(t *testing.T) {
		for _, pClass := range pClasses {
			class, err := svc.Update(context.Background(), pClass.Id, "", time.Now().String())
			assert.NotNil(t, err)
			assert.Nil(t, class)
		}
	})

	t.Run("修改一个不存在的班级", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			class, err := svc.Update(context.Background(), rand.Intn(100)*10000, s, s)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, class)
		}
	})

	t.Run("正常修改", func(t *testing.T) {
		for _, pClass := range pClasses {
			s := time.Now().String()
			name := s + "name"
			desc := s + "desc"
			class, err := svc.Update(context.Background(), pClass.Id, name, desc)
			if assert.Nil(t, err) {
				assert.Equal(t, name, class.Name)
				assert.Equal(t, desc, class.Description)
				assert.Greater(t, class.UpdatedAt.UnixNano(), pClass.UpdatedAt.UnixNano())
			}
		}
	})

	truncateTable(db)
}

func TestClassSvc_Delete(t *testing.T) {
	pClasses, _ := prepareClass(t, db)
	svc := NewClass(classDao)

	t.Run("正常删除", func(t *testing.T) {
		for _, pClass := range pClasses {
			err := svc.Delete(context.Background(), pClass.Id)
			if assert.Nil(t, err) {
				var class model.Class
				err = db.Model(&class).Where("id = ?", pClass.Id).Select()
				assert.Zero(t, class)
				assert.Equal(t, pg.ErrNoRows, err)
			}
		}
	})

	t.Run("删除一个不存在的班级", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			err := svc.Delete(context.Background(), rand.Intn(100)*10000)
			assert.Nil(t, err)
		}
	})

	truncateTable(db)
}

func TestClassSvc_IsNameExist(t *testing.T) {
	pClasses, _ := prepareClass(t, db)
	svc := NewClass(classDao)

	t.Run("排除当前班级后，查找当前班级的名称", func(t *testing.T) {
		for _, pClass := range pClasses {
			is, err := svc.IsNameExist(context.Background(), pClass.Name, pClass.Id)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	t.Run("排除当前班级后，查找其他班级的名称", func(t *testing.T) {
		for i := 0; i < len(pClasses)-1; i++ {
			current := pClasses[i]
			next := pClasses[i+1]
			is, err := svc.IsNameExist(context.Background(), next.Name, current.Id)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找存在的名称", func(t *testing.T) {
		for _, pClass := range pClasses {
			is, err := svc.IsNameExist(context.Background(), pClass.Name, 0)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找不存在的名称", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			is, err := svc.IsNameExist(context.Background(), time.Now().String(), 0)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	truncateTable(db)
}
