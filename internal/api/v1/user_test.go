package v1

import (
	"github.com/go-pg/pg/v10"
	"github.com/kataras/iris/v12/httptest"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/testdb"
	"github.com/xuxusheng/time-frequency-be/internal/model"
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

func TestUser_Create(t *testing.T) {
	// 准备数据
	_ = prepareUser(t, db)
	// 初始化 app 和 e，封装成公共函数
	app := testdb.NewApp()
	userController := NewUser(userSvc)
	app.Post("/api/v1/users", userController.Create)

	t.Run("正常创建", func(t *testing.T) {
		e := httptest.New(t, app, httptest.URL("/api/v1/users"))
		s := time.Now().String()
		name := s + "name"
		phone := s + "phone"
		email := s + "email"
		password := s + "password"
		e.POST("").WithJSON(CreateUserReq{
			Name:     name,
			Phone:    phone,
			Email:    email,
			Password: password,
		}).Expect().Status(httptest.StatusOK)
	})

	t.Run("用户名、手机号、邮箱或密码为空", func(t *testing.T) {
		e := httptest.New(t, app, httptest.URL("/api/v1/users"))
		s := time.Now().String()
		name := s + "name"
		phone := s + "phone"
		email := s + "email"
		password := s + "password"

		e.POST("").WithJSON(CreateUserReq{
			Name:     "",
			Phone:    phone,
			Email:    email,
			Password: password,
		}).Expect().Status(httptest.StatusBadRequest)

		e.POST("").WithJSON(CreateUserReq{
			Name:     name,
			Phone:    "",
			Email:    email,
			Password: password,
		}).Expect().Status(httptest.StatusBadRequest)
	})

	// 清空数据
	_ = testdb.Truncate(db)
}

//func TestUser_Get(t *testing.T) {
//	pUsers := prepareUser(t, db)
//	app := testdb.NewApp()
//	userController := NewUser(userSvc)
//	app.Post("/api/v1/user/get", userController.Get)
//
//	t.Run("正常获取", func(t *testing.T) {
//		e := httptest.New(t, app, httptest.URL("/api/v1/user/get"), httptest.Debug(true), httptest.LogLevel("debug"))
//		for _, pUser := range pUsers {
//
//			e.POST("").WithJSON(UserGetReq{
//				Id: pUser.Id,
//			}).Expect().Status(httptest.StatusOK)
//
//		}
//	})
//
//	// 清空数据
//	_ = testdb.Truncate(db)
//}
