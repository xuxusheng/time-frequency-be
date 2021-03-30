package service

import (
	"context"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
)

type IUser interface {
	Create(ctx context.Context, createdById int, name, phone, email, password string) (*model.User, error)
	Get(ctx context.Context, id int) (*model.User, error)
	GetByName(ctx context.Context, name string) (*model.User, error)
	ListAndCount(ctx context.Context, query string, p *model.Page) ([]*model.User, int, error)
	Update(ctx context.Context, id int, name, phone, email string) (*model.User, error)
	Delete(ctx context.Context, id int) error
	IsNameExist(ctx context.Context, name string, excludeId int) (bool, error)   // 查询用户名是否被占用
	IsPhoneExist(ctx context.Context, phone string, excludeId int) (bool, error) // 查询手机号是否被占用
	IsEmailExist(ctx context.Context, email string, excludeId int) (bool, error) // 查询邮箱是否被占用
}

func NewUser(dao dao.IUser) *User {
	return &User{Dao: dao}
}

type User struct {
	Dao dao.IUser
}

func (u *User) Create(ctx context.Context, createdById int, name, phone, email, password string) (*model.User, error) {
	d := u.Dao
	// 判断用户名是否已存在
	is, err := d.IsNameExist(ctx, name, 0)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, cerror.BadRequest.WithMsg("用户名已存在")
	}

	// 判断手机号是否已存在
	is, err = d.IsPhoneExist(ctx, phone, 0)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, cerror.BadRequest.WithMsg("手机号已存在")
	}

	// 判断邮箱是否已存在
	is, err = d.IsEmailExist(ctx, email, 0)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, cerror.BadRequest.WithMsg("邮箱已存在")
	}

	user, err := d.Create(ctx, createdById, name, phone, email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Get(ctx context.Context, id int) (*model.User, error) {
	return u.Dao.Get(ctx, id)
}

func (u *User) GetByName(ctx context.Context, name string) (*model.User, error) {
	return u.Dao.GetByName(ctx, name)
}

func (u *User) ListAndCount(ctx context.Context, query string, p *model.Page) ([]*model.User, int, error) {
	return u.Dao.ListAndCount(ctx, query, p)
}

func (u *User) Update(ctx context.Context, id int, name, phone, email string) (*model.User, error) {
	d := u.Dao
	// 判断用户名是否被占用
	is, err := d.IsNameExist(ctx, name, id)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, cerror.BadRequest.WithMsg("用户名已存在")
	}
	is, err = d.IsPhoneExist(ctx, phone, id)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, cerror.BadRequest.WithMsg("手机号已存在")
	}
	is, err = d.IsEmailExist(ctx, email, id)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, cerror.BadRequest.WithMsg("邮箱已存在")
	}

	user, err := d.Update(ctx, id, name, phone, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Delete(ctx context.Context, id int) error {
	return u.Dao.Delete(ctx, id)
}

func (u *User) IsNameExist(ctx context.Context, name string, excludeId int) (bool, error) {
	return u.Dao.IsNameExist(ctx, name, excludeId)
}

func (u *User) IsPhoneExist(ctx context.Context, phone string, excludeId int) (bool, error) {
	return u.Dao.IsPhoneExist(ctx, phone, excludeId)
}

func (u *User) IsEmailExist(ctx context.Context, email string, excludeId int) (bool, error) {
	return u.Dao.IsEmailExist(ctx, email, excludeId)
}
