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

// 准备用户数据
func prepareUser(t *testing.T, db *pg.DB) []*model.User {
	users, err := testdb.SeedUser(db)
	if err != nil && len(users) == 0 {
		t.Fatalf("准备用户数据失败：%v", err)
	}
	return users
}

func TestUserDao_Create(t *testing.T) {
	pUsers := prepareUser(t, db)
	dao := NewUserDao(db)

	t.Run("正常创建", func(t *testing.T) {
		at := assert.New(t)
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			user, err := dao.Create(context.Background(), 1, s, s, s, s)
			if at.Nil(err) {
				at.NotZero(user.Id)
				at.Equal(s, user.Name)
				at.Equal(s, user.Phone)
				at.Equal(s, user.Email)
				at.NotNil(s, user.UpdatedAt)
				at.NotNil(s, user.CreatedAt)
			}
		}
	})

	t.Run("用户名、手机号或邮箱重复", func(t *testing.T) {
		for _, pUser := range pUsers {
			s := time.Now().String()
			t.Run("用户名重复", func(t *testing.T) {
				_, err := dao.Create(context.Background(), 1, pUser.Name, s, s, s)
				assert.NotNil(t, err)
			})
			t.Run("手机号重复", func(t *testing.T) {
				_, err := dao.Create(context.Background(), 1, s, pUser.Phone, s, s)
				assert.NotNil(t, err)
			})
			t.Run("邮箱重复", func(t *testing.T) {
				_, err := dao.Create(context.Background(), 1, s, s, pUser.Email, s)
				assert.NotNil(t, err)
			})

		}
	})

	t.Run("用户名、手机号或邮箱为空", func(t *testing.T) {
		// todo
	})

	// 清空数据库
	truncateTable(db)
}

func TestUserDao_Get(t *testing.T) {
	pUsers := prepareUser(t, db)
	dao := NewUserDao(db)

	t.Run("正常获取", func(t *testing.T) {
		for _, pUser := range pUsers {
			user, err := dao.Get(context.Background(), pUser.Id)
			if assert.Nil(t, err) {
				assert.Equal(t, pUser, user)
			}
		}
	})

	t.Run("获取不存在的用户", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			user, err := dao.Get(context.Background(), rand.Intn(100)*10000)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, user)
		}
	})

	// 清空数据库
	truncateTable(db)
}

func TestUserDao_Update(t *testing.T) {
	pUsers := prepareUser(t, db)
	dao := NewUserDao(db)

	t.Run("用户名、手机号或邮箱重复", func(t *testing.T) {
		for i := 0; i < len(pUsers)-1; i++ {
			currentUser := pUsers[i]
			nextUser := pUsers[i+1]
			s := time.Now().String()
			// 用户名重复
			user, err := dao.Update(context.Background(), currentUser.Id, nextUser.Name, s, s)
			assert.Nil(t, user)
			assert.NotNil(t, err)
			// 手机号重复
			user, err = dao.Update(context.Background(), currentUser.Id, s, nextUser.Phone, s)
			assert.Nil(t, user)
			assert.NotNil(t, err)
			// 邮箱重复
			user, err = dao.Update(context.Background(), currentUser.Id, s, s, nextUser.Email)
			assert.Nil(t, user)
			assert.NotNil(t, err)
		}
	})

	t.Run("用户名、手机号或邮箱为空", func(t *testing.T) {
		for _, pUser := range pUsers {
			s := time.Now().String()
			// 用户名为空
			user, err := dao.Update(context.Background(), pUser.Id, "", s, s)
			assert.Nil(t, user)
			assert.NotNil(t, err)
			// 手机号为空
			user, err = dao.Update(context.Background(), pUser.Id, s, "", s)
			assert.Nil(t, user)
			assert.NotNil(t, err)
			// 邮箱为空
			user, err = dao.Update(context.Background(), pUser.Id, s, s, "")
			assert.Nil(t, user)
			assert.NotNil(t, err)
		}
	})

	t.Run("修改一个不存在的用户", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			user, err := dao.Update(context.Background(), rand.Intn(100)*10000, "", "", "")
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, user)
		}
	})

	t.Run("正常修改", func(t *testing.T) {
		// 等待一会儿，让 updatedAt 字段变化
		for _, pUser := range pUsers {
			// 将时间字符串，用作修改的值，避免重复
			s := time.Now().String()
			name := s + "name"
			phone := s + "phone"
			email := s + "email"
			user, err := dao.Update(context.Background(), pUser.Id, name, phone, email)
			if assert.Nil(t, err) {
				assert.Equal(t, name, user.Name)
				assert.Equal(t, phone, user.Phone)
				assert.Equal(t, email, user.Email)
				assert.Greater(t, user.UpdatedAt.UnixNano(), pUser.UpdatedAt.UnixNano())
			}
		}
	})

	// 清空数据库
	truncateTable(db)
}

func TestUserDao_Delete(t *testing.T) {

	pUsers := prepareUser(t, db)
	dao := NewUserDao(db)

	t.Run("正常删除", func(t *testing.T) {
		at := assert.New(t)
		for _, pUser := range pUsers {
			err := dao.Delete(context.Background(), pUser.Id)
			if at.Nil(err) {
				var user model.User
				err = db.Model(&user).Where("id = ?", pUser.Id).Select()
				at.Zero(user)
				at.Equal(pg.ErrNoRows, err)
			}
		}
	})

	t.Run("删除不存在的用户", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			err := dao.Delete(context.Background(), rand.Intn(100)*10000)
			assert.Nil(t, err)
		}
	})
}

func TestUserDao_IsNameExist(t *testing.T) {
	pUsers := prepareUser(t, db)
	dao := NewUserDao(db)

	t.Run("排除当前用户之后，查找当前用户的用户名", func(t *testing.T) {
		for _, pUser := range pUsers {
			is, err := dao.IsNameExist(context.Background(), pUser.Name, pUser.Id)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	t.Run("排除当前用户之后，查找其他用户的用户名", func(t *testing.T) {
		for i := 0; i < len(pUsers)-1; i++ {
			current := pUsers[i]
			next := pUsers[i+1]
			is, err := dao.IsNameExist(context.Background(), next.Name, current.Id)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("不排除当前用户，查找用户名", func(t *testing.T) {
		for _, pUser := range pUsers {
			is, err := dao.IsNameExist(context.Background(), pUser.Name, 0)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找不存在的用户名", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			is, err := dao.IsNameExist(context.Background(), time.Now().String(), 0)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	truncateTable(db)
}

func TestUserDao_IsPhoneExist(t *testing.T) {
	pUsers := prepareUser(t, db)
	dao := NewUserDao(db)

	t.Run("排除当前用户之后，查找当前用户的手机号", func(t *testing.T) {
		for _, pUser := range pUsers {
			is, err := dao.IsPhoneExist(context.Background(), pUser.Phone, pUser.Id)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	t.Run("排除当前用户之后，查找其他用户的手机号", func(t *testing.T) {
		for i := 0; i < len(pUsers)-1; i++ {
			current := pUsers[i]
			next := pUsers[i+1]
			is, err := dao.IsPhoneExist(context.Background(), next.Phone, current.Id)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("不排除当前用户，查找手机号", func(t *testing.T) {
		for _, pUser := range pUsers {
			is, err := dao.IsPhoneExist(context.Background(), pUser.Phone, 0)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找不存在的手机号", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			is, err := dao.IsNameExist(context.Background(), time.Now().String(), 0)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	truncateTable(db)
}

func TestUserDao_IsEmailExist(t *testing.T) {
	pUsers := prepareUser(t, db)
	dao := NewUserDao(db)

	t.Run("排除当前用户之后，查找当前用户的邮箱", func(t *testing.T) {
		for _, pUser := range pUsers {
			is, err := dao.IsEmailExist(context.Background(), pUser.Email, pUser.Id)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	t.Run("排除当前用户之后，查找其他用户的邮箱", func(t *testing.T) {
		for i := 0; i < len(pUsers)-1; i++ {
			current := pUsers[i]
			next := pUsers[i+1]
			is, err := dao.IsEmailExist(context.Background(), next.Email, current.Id)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("不排除当前用户，查找邮箱", func(t *testing.T) {
		for _, pUser := range pUsers {
			is, err := dao.IsEmailExist(context.Background(), pUser.Email, 0)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找不存在的邮箱", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			is, err := dao.IsNameExist(context.Background(), time.Now().String(), 0)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	truncateTable(db)
}
