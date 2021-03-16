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

// 准备班级数据
func prepareClass(t *testing.T, db *pg.DB) ([]*model.Class, []*model.User) {
	users, err := testdb.SeedUser(db)
	if err != nil {
		t.Fatalf("准备用户数据失败：%v", err)
	}
	classes, err := testdb.SeedClass(db, users)
	if err != nil && len(classes) == 0 {
		t.Fatalf("准备班级数据失败：%v", err)
	}
	return classes, users
}

func TestClassDao_Create(t *testing.T) {
	pClasses, pUsers := prepareClass(t, db)
	dao := NewClassDao(db)

	createdById := pUsers[0].Id

	t.Run("班级名称重复", func(t *testing.T) {
		for _, pClass := range pClasses {
			class, err := dao.Create(context.Background(), createdById, pClass.Name, time.Now().String())
			assert.NotNil(t, err)
			assert.Nil(t, class)
		}
	})

	t.Run("班级名称为空", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			class, err := dao.Create(context.Background(), createdById, "", time.Now().String())
			assert.NotNil(t, err)
			assert.Nil(t, class)
		}
	})

	t.Run("描述为空", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			class, err := dao.Create(context.Background(), createdById, time.Now().String(), "")
			assert.Nil(t, err)
			assert.NotZero(t, class)
		}
	})

	t.Run("正常创建", func(t *testing.T) {
		at := assert.New(t)
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			name := s + "name"
			description := s + "description"
			class, err := dao.Create(context.Background(), createdById, name, description)
			if at.Nil(err) {
				at.NotZero(class.Id)
				at.Equal(name, class.Name)
				at.Equal(description, class.Description)
				at.NotNil(s, class.CreatedAt)
				at.NotNil(s, class.UpdatedAt)
			}
		}
	})

	// 清空数据库
	truncateTable(db)
}

func TestClassDao_Get(t *testing.T) {
	pClasses, _ := prepareClass(t, db)
	dao := NewClassDao(db)

	t.Run("正常获取", func(t *testing.T) {
		for _, pClass := range pClasses {
			class, err := dao.Get(context.Background(), pClass.Id)
			if assert.Nil(t, err) {
				assert.Equal(t, pClass, class)
			}
		}
	})

	t.Run("获取不存在的班级", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			class, err := dao.Get(context.Background(), rand.Intn(100)*10000)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, class)
		}
	})

	truncateTable(db)
}

func TestClassDao_Update(t *testing.T) {
	pClasses, _ := prepareClass(t, db)
	dao := NewClassDao(db)

	t.Run("班级名称重复", func(t *testing.T) {
		for i := 0; i < len(pClasses)-1; i++ {
			current := pClasses[i]
			next := pClasses[i+1]
			class, err := dao.Update(context.Background(), current.Id, next.Name, time.Now().String())
			assert.Nil(t, class)
			assert.NotNil(t, err)
		}
	})

	t.Run("班级名称为空", func(t *testing.T) {
		for _, pClass := range pClasses {
			class, err := dao.Update(context.Background(), pClass.Id, "", time.Now().String())
			assert.NotNil(t, err)
			assert.Nil(t, class)
		}
	})

	t.Run("描述为空", func(t *testing.T) {
		for _, pClass := range pClasses {
			class, err := dao.Update(context.Background(), pClass.Id, time.Now().String(), "")
			assert.Nil(t, err)
			assert.Zero(t, class.Description)
		}
	})

	t.Run("修改一个不存在的班级", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			class, err := dao.Update(context.Background(), rand.Intn(100)*10000, s, s)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, class)
		}
	})

	t.Run("正常修改", func(t *testing.T) {
		// 等待一会儿，让 updatedAt 字段变化
		for _, pClass := range pClasses {
			s := time.Now().String()
			name := s + "name"
			description := s + "desc"
			class, err := dao.Update(context.Background(), pClass.Id, name, description)
			if assert.Nil(t, err) {
				assert.Equal(t, name, class.Name)
				assert.Equal(t, description, class.Description)
				assert.Greater(t, class.UpdatedAt.UnixNano(), pClass.UpdatedAt.UnixNano())
			}
		}
	})

	// 清空数据库
	truncateTable(db)
}

func TestClassDao_Delete(t *testing.T) {
	pClasses, _ := prepareClass(t, db)
	dao := NewClassDao(db)

	t.Run("正常删除", func(t *testing.T) {
		at := assert.New(t)
		for _, pClass := range pClasses {
			err := dao.Delete(context.Background(), pClass.Id)
			if at.Nil(err) {

				var class model.Class
				err = db.Model(&class).Where("id = ?", pClass.Id).Select()
				at.Zero(class)
				at.Equal(pg.ErrNoRows, err)

			}
		}
	})

	t.Run("删除不存在的班级", func(t *testing.T) {
		for i := 0; i < 0; i++ {
			err := dao.Delete(context.Background(), rand.Intn(100)*10000)
			assert.Nil(t, err)
		}
	})

	truncateTable(db)
}
