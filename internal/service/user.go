package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/pkg/errcode"
	"gorm.io/gorm"
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
func (u *UserService) Create(name, phone, password string) (*model.User, *errcode.Error) {

	userDao := dao.NewUserDao(global.DBEngine)

	// 检查用户名
	isExist, err := userDao.IsNameExist(name)
	if err != nil {
		return nil, errcode.CreateUserFail.WithDetails(err.Error())
	}
	if isExist {
		return nil, errcode.CreateUserFailNameExist
	}

	// 检查手机号
	isExist, err = userDao.IsPhoneExist(phone)
	if err != nil {
		return nil, errcode.CreateUserFail.WithDetails(err.Error())
	}
	if isExist {
		return nil, errcode.CreateUserFailPhoneExist
	}

	// 检查无误，开始写入
	user, err := userDao.Create(name, phone, password)
	if err != nil {
		return nil, errcode.CreateUserFail.WithDetails(err.Error())
	}

	return user, nil
}

func (u *UserService) Delete(id uint) error {
	userDao := dao.NewUserDao(global.DBEngine)
	return userDao.Delete(id)
}

// 更新用户基本信息，在 server 外层，就应该先校验好 name 和 phone 不同时为空字符串
func (u *UserService) Update(id uint, name, phone string) *errcode.Error {
	userDao := dao.NewUserDao(global.DBEngine)

	if name != "" {
		// 检查用户名
		isExist, err := userDao.IsNameExist(name)
		if err != nil {
			return errcode.UpdateUserFail.WithDetails(err.Error())
		}
		if isExist {
			return errcode.UpdateUserFailNameExist
		}
	}

	if phone != "" {
		// 检查手机号
		isExist, err := userDao.IsNameExist(name)
		if err != nil {
			return errcode.UpdateUserFail.WithDetails(err.Error())
		}
		if isExist {
			return errcode.UpdateUserFailPhoneExist
		}
	}

	err := userDao.Update(id, gin.H{
		"name":  name,
		"phone": phone,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.NotFound.WithMsg("用户不存在")
		}
		return errcode.UpdateUserFail.WithDetails(err.Error())
	}
	return nil
}

func (u *UserService) Get(id uint) (*model.User, error) {
	userDao := dao.NewUserDao(global.DBEngine)
	return userDao.Get(id)
}

func (u *UserService) List(name, phone string, pn, ps int) ([]*model.User, int64, error) {
	userDao := dao.NewUserDao(global.DBEngine)

	// 查询命中记录条数
	count, err := userDao.Count(name, phone)
	if err != nil {
		return nil, 0, err
	}

	// 获取当前页数据
	users, err := userDao.List(name, phone, pn, ps)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}
