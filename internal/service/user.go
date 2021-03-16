package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/model"
)

type IUserSvc interface {
	Create(ctx context.Context, createdById int, name, phone, email, password string) (*model.User, error)
	Get(ctx context.Context, id int) (*model.User, error)
	Update(ctx context.Context, id int, name, phone, email string) (*model.User, error)
	Delete(ctx context.Context, id int) error
	IsNameExist(ctx context.Context, name string, excludeId int) (bool, error)   // 查询用户名是否被占用
	IsPhoneExist(ctx context.Context, phone string, excludeId int) (bool, error) // 查询手机号是否被占用
	IsEmailExist(ctx context.Context, email string, excludeId int) (bool, error) // 查询邮箱是否被占用
}

func NewUserSvc(dao dao.IUserDao) *UserSvc {
	return &UserSvc{Dao: dao}
}

type UserSvc struct {
	Dao dao.IUserDao
}

func (u *UserSvc) Create(ctx context.Context, createdById int, name, phone, email, password string) (*model.User, error) {
	d := u.Dao
	// 判断用户名是否已存在
	is, err := d.IsNameExist(ctx, name, 0)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, errors.New("用户名已存在")
	}

	// 判断手机号是否已存在
	is, err = d.IsPhoneExist(ctx, phone, 0)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, errors.New("手机号已存在")
	}

	// 判断邮箱是否已存在
	is, err = d.IsEmailExist(ctx, email, 0)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, errors.New("邮箱已存在")
	}

	user, err := d.Create(ctx, createdById, name, phone, email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserSvc) Get(ctx context.Context, id int) (*model.User, error) {
	return u.Dao.Get(ctx, id)
}

func (u *UserSvc) Update(ctx context.Context, id int, name, phone, email string) (*model.User, error) {
	d := u.Dao
	// 判断用户名是否被占用
	is, err := d.IsNameExist(ctx, name, id)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, errors.New("用户名已存在")
	}
	is, err = d.IsPhoneExist(ctx, phone, id)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, errors.New("手机号已存在")
	}
	is, err = d.IsEmailExist(ctx, email, id)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, errors.New("邮箱已存在")
	}

	user, err := d.Update(ctx, id, name, phone, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserSvc) Delete(ctx context.Context, id int) error {
	return u.Dao.Delete(ctx, id)
}

func (u *UserSvc) IsNameExist(ctx context.Context, name string, excludeId int) (bool, error) {
	return u.Dao.IsNameExist(ctx, name, excludeId)
}

func (u *UserSvc) IsPhoneExist(ctx context.Context, phone string, excludeId int) (bool, error) {
	return u.Dao.IsPhoneExist(ctx, phone, excludeId)
}

func (u *UserSvc) IsEmailExist(ctx context.Context, email string, excludeId int) (bool, error) {
	return u.Dao.IsEmailExist(ctx, email, excludeId)
}
