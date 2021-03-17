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

func prepareUser(t *testing.T, db *pg.DB) []*model.User {
	users, err := testdb.SeedUser(db)
	if err != nil || len(users) == 0 {
		t.Fatalf("准备用户数据失败: %v", err)
	}
	return users
}

func TestUserSvc_Create(t *testing.T) {

	pUsers := prepareUser(t, db)
	svc := NewUserSvc(userDao)

	t.Run("用户名、手机号或邮箱重复", func(t *testing.T) {
		for _, pUser := range pUsers {
			s := time.Now().String()
			name := s + "name"
			phone := s + "phone"
			email := s + "email"
			pwd := s + "password"
			t.Run("用户名重复", func(t *testing.T) {
				_, err := svc.Create(context.Background(), 1, pUser.Name, phone, email, pwd)
				assert.EqualError(t, err, "用户名已存在")
			})
			t.Run("手机号重复", func(t *testing.T) {
				_, err := svc.Create(context.Background(), 1, name, pUser.Phone, email, pwd)
				assert.EqualError(t, err, "手机号已存在")
			})
			t.Run("邮箱重复", func(t *testing.T) {
				_, err := svc.Create(context.Background(), 1, name, phone, pUser.Email, pwd)
				assert.EqualError(t, err, "邮箱已存在")
			})
		}
	})

	t.Run("用户名、手机号、邮箱或密码为空", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			name := s + "name"
			phone := s + "phone"
			email := s + "email"
			pwd := s + "password"
			t.Run("用户名为空", func(t *testing.T) {
				_, err := svc.Create(context.Background(), 1, "", phone, email, pwd)
				assert.NotNil(t, err)
			})
			t.Run("手机号为空", func(t *testing.T) {
				_, err := svc.Create(context.Background(), 1, name, "", email, pwd)
				assert.NotNil(t, err)
			})
			t.Run("邮箱为空", func(t *testing.T) {
				_, err := svc.Create(context.Background(), 1, name, phone, "", pwd)
				assert.NotNil(t, err)
			})
			t.Run("密码为空", func(t *testing.T) {
				_, err := svc.Create(context.Background(), 1, name, phone, email, "")
				assert.NotNil(t, err)
			})
		}
	})

	t.Run("正常创建", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			name := s + "name"
			phone := s + "phone"
			email := s + "email"
			pwd := s + "password"
			user, err := svc.Create(context.Background(), 1, name, phone, email, pwd)
			if assert.Nil(t, err) {
				assert.Equal(t, name, user.Name)
				assert.Equal(t, phone, user.Phone)
				assert.Equal(t, email, user.Email)
				assert.Equal(t, pwd, user.Password)
			}
		}
	})

	truncateTable(db)
}

func TestUserSvc_Get(t *testing.T) {

	pUsers := prepareUser(t, db)
	svc := NewUserSvc(userDao)

	t.Run("正常获取", func(t *testing.T) {
		for _, pUser := range pUsers {
			user, err := svc.Get(context.Background(), pUser.Id)
			if assert.Nil(t, err) {
				assert.Equal(t, pUser, user)
			}
		}
	})

	t.Run("获取不存在的用户", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			user, err := svc.Get(context.Background(), rand.Intn(100)*10000)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, user)
		}
	})

	truncateTable(db)
}

func TestUserSvc_Update(t *testing.T) {
	pUsers := prepareUser(t, db)
	svc := NewUserSvc(userDao)
	t.Run("用户名、手机号或邮箱重复", func(t *testing.T) {
		for i, pUser := range pUsers {
			if i == len(pUsers)-1 {
				break
			}
			next := pUsers[i+1]
			s := time.Now().String()
			name := s + "name"
			phone := s + "phone"
			email := s + "email"
			t.Run("用户名重复", func(t *testing.T) {
				user, err := svc.Update(context.Background(), pUser.Id, next.Name, phone, email)
				assert.EqualError(t, err, "用户名已存在")
				assert.Nil(t, user)
			})
			t.Run("手机号重复", func(t *testing.T) {
				user, err := svc.Update(context.Background(), pUser.Id, name, next.Phone, email)
				assert.EqualError(t, err, "手机号已存在")
				assert.Nil(t, user)
			})
			t.Run("邮箱重复", func(t *testing.T) {
				user, err := svc.Update(context.Background(), pUser.Id, name, phone, next.Email)
				assert.EqualError(t, err, "邮箱已存在")
				assert.Nil(t, user)
			})
		}
	})

	t.Run("用户名、手机号、邮箱为空", func(t *testing.T) {
		for _, pUser := range pUsers {
			s := time.Now().String()
			id := pUser.Id
			name := s + "name"
			phone := s + "phone"
			email := s + "email"
			t.Run("用户名为空", func(t *testing.T) {
				user, err := svc.Update(context.Background(), id, "", phone, email)
				assert.NotNil(t, err)
				assert.Nil(t, user)
			})
			t.Run("手机号为空", func(t *testing.T) {
				user, err := svc.Update(context.Background(), id, name, "", email)
				assert.NotNil(t, err)
				assert.Nil(t, user)
			})
			t.Run("邮箱为空", func(t *testing.T) {
				user, err := svc.Update(context.Background(), id, name, phone, "")
				assert.NotNil(t, err)
				assert.Nil(t, user)
			})
		}
	})

	t.Run("修改一个不存在的用户", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			s := time.Now().String()
			user, err := svc.Update(context.Background(), rand.Intn(100)*10000, s, s, s)
			assert.Equal(t, pg.ErrNoRows, err)
			assert.Nil(t, user)
		}
	})

	t.Run("正常修改", func(t *testing.T) {
		for _, pUser := range pUsers {
			s := time.Now().String()
			name := s + "name"
			phone := s + "phone"
			email := s + "email"
			user, err := svc.Update(context.Background(), pUser.Id, name, phone, email)
			if assert.Nil(t, err) {
				assert.Equal(t, name, user.Name)
				assert.Equal(t, phone, user.Phone)
				assert.Equal(t, email, user.Email)
				assert.Greater(t, user.UpdatedAt.UnixNano(), pUser.UpdatedAt.UnixNano())
			}
		}
	})

	truncateTable(db)
}

func TestUserSvc_Delete(t *testing.T) {
	pUsers := prepareUser(t, db)
	svc := NewUserSvc(userDao)

	t.Run("正常删除", func(t *testing.T) {
		for _, pUser := range pUsers {
			err := svc.Delete(context.Background(), pUser.Id)
			if assert.Nil(t, err) {
				var user model.User
				err = db.Model(&user).Where("id = ?", pUser.Id).Select()
				assert.Zero(t, user)
				assert.Equal(t, pg.ErrNoRows, err)
			}
		}
	})

	t.Run("删除不存在的用户", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			err := svc.Delete(context.Background(), rand.Intn(100)*10000)
			assert.Nil(t, err)
		}
	})
}

func TestUserSvc_IsNameExist(t *testing.T) {
	pUsers := prepareUser(t, db)
	svc := NewUserSvc(userDao)

	t.Run("排除当前用户之后，查找当前用户的用户名", func(t *testing.T) {
		for _, pUser := range pUsers {
			is, err := svc.IsNameExist(context.Background(), pUser.Name, pUser.Id)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	t.Run("排除当前用户之后，查找其他用户的用户名", func(t *testing.T) {
		for i := 0; i < len(pUsers)-1; i++ {
			current := pUsers[i]
			next := pUsers[i+1]
			is, err := svc.IsNameExist(context.Background(), next.Name, current.Id)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找存在的用户名", func(t *testing.T) {
		for _, pUser := range pUsers {
			is, err := svc.IsNameExist(context.Background(), pUser.Name, 0)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找不存在的用户名", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			is, err := svc.IsNameExist(context.Background(), time.Now().String(), 0)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	truncateTable(db)
}

func TestUserSvc_IsPhoneExist(t *testing.T) {
	pUsers := prepareUser(t, db)
	svc := NewUserSvc(userDao)

	t.Run("排除当前用户之后，查找当前用户的手机号", func(t *testing.T) {
		for _, pUser := range pUsers {
			is, err := svc.IsPhoneExist(context.Background(), pUser.Phone, pUser.Id)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	t.Run("排除当前用户之后，查找其他用户的手机号", func(t *testing.T) {
		for i := 0; i < len(pUsers)-1; i++ {
			current := pUsers[i]
			next := pUsers[i+1]
			is, err := svc.IsPhoneExist(context.Background(), next.Phone, current.Id)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("不排除当前用户，查找手机号", func(t *testing.T) {
		for _, pUser := range pUsers {
			is, err := svc.IsPhoneExist(context.Background(), pUser.Phone, 0)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找不存在的手机号", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			is, err := svc.IsNameExist(context.Background(), time.Now().String(), 0)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	truncateTable(db)
}

func TestUserSvc_IsEmailExist(t *testing.T) {
	pUsers := prepareUser(t, db)
	svc := NewUserSvc(userDao)

	t.Run("排除当前用户之后，查找当前用户的邮箱", func(t *testing.T) {
		for _, pUser := range pUsers {
			is, err := svc.IsEmailExist(context.Background(), pUser.Email, pUser.Id)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	t.Run("排除当前用户之后，查找其他用户的邮箱", func(t *testing.T) {
		for i := 0; i < len(pUsers)-1; i++ {
			current := pUsers[i]
			next := pUsers[i+1]
			is, err := svc.IsEmailExist(context.Background(), next.Email, current.Id)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找存在的邮箱", func(t *testing.T) {
		for _, pUser := range pUsers {
			is, err := svc.IsEmailExist(context.Background(), pUser.Email, 0)
			if assert.Nil(t, err) {
				assert.True(t, is)
			}
		}
	})

	t.Run("查找不存在的邮箱", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			is, err := svc.IsNameExist(context.Background(), time.Now().String(), 0)
			if assert.Nil(t, err) {
				assert.False(t, is)
			}
		}
	})

	truncateTable(db)
}
