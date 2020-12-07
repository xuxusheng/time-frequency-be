package service

import (
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/pkg/app"
	"gorm.io/gorm"
)

type CountUserReq struct {
	Name  string
	Phone string
}

func (svc *Service) CountUser(param *CountUserReq) (int64, error) {
	return svc.dao.CountUser(param.Name, param.Phone)
}

type UserListReq struct {
	Name  string `form:"name" binding:"max=100"`
	Phone string `form:"phone" binding:"max=20"`
	Pn    string `form:"pn"`
	Ps    string `form:"ps"`
}

func (svc *Service) GetUserList(param *UserListReq, pageInfo *app.PageInfo) ([]*model.User, error) {
	return svc.dao.GetUserList(
		param.Name,
		param.Phone,
		pageInfo.Pn,
		pageInfo.Ps,
	)
}

type CreateUserReq struct {
	Name     string `form:"name" binding:"min=1"`
	Phone    string `form:"phone" binding:"min=1"`
	Password string `form:"password" binding:"min=1"`
}

func (svc *Service) CreateUser(param *CreateUserReq) error {

	// 针对密码进行加密

	return svc.dao.CreateUser(
		param.Name,
		param.Phone,
		param.Password,
	)
}

// 用来判断用户名和手机号是否已被占用
func (svc *Service) IsUserExist(name, phone string) (bool, error) {
	_, err := svc.dao.GetUser(0, name, phone)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 错误为 notFount，不存在此记录
			return false, nil
		}
		// 其他错误，查询失败
		return false, err
	}
	// 查询未出错，存在此条记录
	return true, nil
}
