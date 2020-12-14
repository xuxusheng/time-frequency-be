package service

import (
	"context"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/pg_dao"
	"github.com/xuxusheng/time-frequency-be/internal/pg_model"
	"github.com/xuxusheng/time-frequency-be/pkg/app"
	"github.com/xuxusheng/time-frequency-be/pkg/errcode"
)

type UserService struct {
	Service
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{
		Service: Service{ctx: ctx},
	}
}

// 创建用户，需要区分不能的错误类型，所有在 service 层中直接抛出 errcode.Error 类型，如果不需要区分的话，可以直接抛标准 Error 类型就行
func (u *UserService) Create(name, phone, password string) (*pg_model.User, *errcode.Error) {

	//userDao := dao.NewUserDao(global.DBEngine)
	userDao := pg_dao.NewUserDao(global.PGEngine)

	// 检查用户名
	isExist, err := userDao.IsNameExist(name, 0)
	if err != nil {
		return nil, errcode.CreateUserFail.WithDetails(err.Error())
	}
	if isExist {
		return nil, errcode.CreateUserFailNameExist
	}

	// 检查手机号
	isExist, err = userDao.IsPhoneExist(phone, 0)
	if err != nil {
		return nil, errcode.CreateUserFail.WithDetails(err.Error())
	}
	if isExist {
		return nil, errcode.CreateUserFailPhoneExist
	}

	// 密码算 hash
	hash, err := app.EncodePWD(password)
	if err != nil {
		return nil, errcode.CreateUserFail.WithDetails(err.Error())
	}

	// 检查无误，开始写入
	user, err := userDao.Create(name, phone, hash)
	if err != nil {
		return nil, errcode.CreateUserFail.WithDetails(err.Error())
	}

	return user, nil
}

func (u *UserService) Delete(id int) *errcode.Error {
	//userDao := dao.NewUserDao(global.DBEngine)
	userDao := pg_dao.NewUserDao(global.PGEngine)

	// 判断用户是否存在
	isExist, err := userDao.IsIDExist(id)
	if err != nil {
		return errcode.DeleteUserFail.WithDetails(err.Error())
	}
	if !isExist {
		return errcode.NotFound.WithMsg("用户不存在")
	}

	// 执行删除
	if err := userDao.Delete(id); err != nil {
		return errcode.DeleteUserFail.WithDetails(err.Error())
	}
	return nil
}

// 更新用户基本信息，在 server 外层，就应该先校验好 name 和 phone 不同时为空字符串
func (u *UserService) Update(id int, name, phone string) *errcode.Error {
	//userDao := dao.NewUserDao(global.DBEngine)
	userDao := pg_dao.NewUserDao(global.PGEngine)

	if name != "" {
		// 检查用户名
		isExist, err := userDao.IsNameExist(name, id)
		if err != nil {
			return errcode.UpdateUserFail.WithDetails(err.Error())
		}
		if isExist {
			return errcode.UpdateUserFailNameExist
		}
	}

	if phone != "" {
		// 检查手机号
		isExist, err := userDao.IsPhoneExist(phone, id)
		if err != nil {
			return errcode.UpdateUserFail.WithDetails(err.Error())
		}
		if isExist {
			return errcode.UpdateUserFailPhoneExist
		}
	}

	// 判断用户是否存在
	isExist, err := userDao.IsIDExist(id)
	if err != nil {
		return errcode.UpdateUserFail.WithDetails(err.Error())
	}
	if !isExist {
		return errcode.NotFound.WithMsg("用户不存在")
	}

	//err = userDao.Update(id, map[string]interface{}{
	//	"name":  name,
	//	"phone": phone,
	//})
	err = userDao.Update(id, name, phone)
	if err != nil {
		return errcode.UpdateUserFail.WithDetails(err.Error())
	}
	return nil
}

func (u *UserService) Get(id int) (*pg_model.User, error) {
	//userDao := dao.NewUserDao(global.DBEngine)
	userDao := pg_dao.NewUserDao(global.PGEngine)

	return userDao.Get(id)
}

func (u *UserService) List(name, phone string, page *model.Page) ([]*pg_model.User, int, error) {
	//userDao := dao.NewUserDao(global.DBEngine)

	//// 查询命中记录条数
	//count, err := userDao.Count(name, phone)
	//if err != nil {
	//	return nil, 0, err
	//}
	//
	//// 获取当前页数据
	//users, err := userDao.List(name, phone, pn, ps)
	//if err != nil {
	//	return nil, 0, err
	//}
	//
	//return users, count, nil

	userDao := pg_dao.NewUserDao(global.PGEngine)

	users, count, err := userDao.ListAndCount(name, phone, page)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil

}
