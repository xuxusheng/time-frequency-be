package service

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/model"
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

// 创建用户，需要区分不能的错误类型，所有在 service 层中直接抛出 errcode.Error 类型，如果不需要区分的话，可以直接抛标准 ToError 类型就行
func (u *UserService) Create(name, phone, password string) (*model.User, *errcode.Error) {

	//userDao := dao.NewUserDao(global.DBEngine)
	userDao := dao.NewUserDao(global.PGEngine)

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

func (u *UserService) Delete(id uint) *errcode.Error {
	//userDao := dao.NewUserDao(global.DBEngine)
	userDao := dao.NewUserDao(global.PGEngine)

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
func (u *UserService) Update(id uint, name, phone string) (*model.User, *errcode.Error) {
	//userDao := dao.NewUserDao(global.DBEngine)
	userDao := dao.NewUserDao(global.PGEngine)

	if name != "" {
		// 检查用户名
		isExist, err := userDao.IsNameExist(name, id)
		if err != nil {
			return nil, errcode.UpdateUserFail.WithDetails(err.Error())
		}
		if isExist {
			return nil, errcode.UpdateUserFailNameExist
		}
	}

	if phone != "" {
		// 检查手机号
		isExist, err := userDao.IsPhoneExist(phone, id)
		if err != nil {
			return nil, errcode.UpdateUserFail.WithDetails(err.Error())
		}
		if isExist {
			return nil, errcode.UpdateUserFailPhoneExist
		}
	}

	// 判断用户是否存在
	isExist, err := userDao.IsIDExist(id)
	if err != nil {
		return nil, errcode.UpdateUserFail.WithDetails(err.Error())
	}
	if !isExist {
		return nil, errcode.NotFound.WithMsg("用户不存在")
	}

	user, err := userDao.Update(id, name, phone)
	if err != nil {
		return nil, errcode.UpdateUserFail.WithDetails(err.Error())
	}
	return user, nil
}

func (u *UserService) UpdatePassword(id uint, oldPassword, newPassword string) *errcode.Error {

	userDao := dao.NewUserDao(global.PGEngine)

	// 判断旧密码是否正确
	user, err := userDao.Get(id)
	if err != nil {
		return errcode.InternalServerError.WithDetails(err.Error())
	}

	err = app.ComparePWD(user.Password, oldPassword)
	if err != nil {
		return errcode.BadRequest.WithMsg("旧密码不正确")
	}

	// 验证通过，开始修改
	hash, err := app.EncodePWD(newPassword)
	if err != nil {
		return errcode.InternalServerError.WithDetails(err.Error())
	}

	err = userDao.UpdatePassword(id, hash)
	if err != nil {
		return errcode.InternalServerError.WithDetails(err.Error())
	}
	return nil
}

func (u *UserService) UpdateRole(id uint, role model.Role) error {
	userDao := dao.NewUserDao(global.PGEngine)
	err := userDao.UpdateRole(id, role)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) Get(id uint) (*model.User, error) {
	userDao := dao.NewUserDao(global.PGEngine)

	return userDao.Get(id)
}

func (u *UserService) List(name, phone string, page *model.Page) ([]*model.User, int, error) {
	userDao := dao.NewUserDao(global.PGEngine)

	users, count, err := userDao.ListAndCount(name, phone, page)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil

}

// 校验密码
func (u *UserService) Login(name, password string) (string, *errcode.Error) {
	userDao := dao.NewUserDao(global.PGEngine)

	// 从数据库中取出 user
	user, err := userDao.GetByName(name)
	if err != nil {
		if err == pg.ErrNoRows {
			return "", errcode.UnauthorizedUserError
		}
		return "", errcode.InternalServerError.WithDetails(err.Error())
	}

	// 校验密码
	err = app.ComparePWD(user.Password, password)
	if err != nil {
		return "", errcode.UnauthorizedUserError
	}

	// 校验通过，生成 token
	signer := jwt.NewSigner(
		jwt.HS256,
		global.JWTSetting.Secret,
		global.JWTSetting.Expire,
	)
	claims := model.JWTClaims{
		UID:   user.ID,
		Roles: []model.Role{model.Role(user.Role)},
	}
	token, err := signer.Sign(claims)
	if err != nil {
		return "", errcode.InternalServerError.WithDetails(err.Error())
	}

	return string(token), nil
}
