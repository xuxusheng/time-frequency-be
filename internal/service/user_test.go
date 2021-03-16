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
	if err != nil && len(users) == 0 {
		t.Fatalf("准备用户数据失败: %v", err)
	}
	return users
}

func TestUserSvc_Create(t *testing.T) {

	pUsers := prepareUser(t, db)
	userSvc := NewUserSvc(userDao)

	t.Run("用户名、手机号或邮箱重复", func(t *testing.T) {
		for _, pUser := range pUsers {
			t.Run("用户名重复", func(t *testing.T) {
				s := time.Now().String()
				_, err := userSvc.Create(context.Background(), 1, pUser.Name, s, s, s)
				assert.EqualError(t, err, "用户名已存在")
			})
			t.Run("手机号重复", func(t *testing.T) {
				s := time.Now().String()
				_, err := userSvc.Create(context.Background(), 1, s, pUser.Phone, s, s)
				assert.EqualError(t, err, "手机号已存在")
			})
			t.Run("邮箱重复", func(t *testing.T) {
				s := time.Now().String()
				_, err := userSvc.Create(context.Background(), 1, s, s, pUser.Email, s)
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
			password := s + "password"
			t.Run("用户名为空", func(t *testing.T) {
				_, err := userSvc.Create(context.Background(), 1, "", phone, email, password)
				assert.NotNil(t, err)
			})
			t.Run("手机号为空", func(t *testing.T) {
				_, err := userSvc.Create(context.Background(), 1, name, "", email, password)
				assert.NotNil(t, err)
			})
			t.Run("邮箱为空", func(t *testing.T) {
				_, err := userSvc.Create(context.Background(), 1, name, phone, "", password)
				assert.NotNil(t, err)
			})
			t.Run("密码为空", func(t *testing.T) {
				_, err := userSvc.Create(context.Background(), 1, name, phone, email, "")
				assert.NotNil(t, err)
			})
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
	//pUsers := prepareUser(t, db)
	//svc := NewUserSvc(userDao)

	t.Run("用户名、手机号或邮箱重复", func(t *testing.T) {

	})

	t.Run("用户名、手机号或邮箱为空", func(t *testing.T) {

	})

	t.Run("修改一个不存在的用户", func(t *testing.T) {

	})

	t.Run("正常修改", func(t *testing.T) {

	})

	truncateTable(db)
}
