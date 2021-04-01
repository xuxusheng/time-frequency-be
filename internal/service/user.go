package service

import (
	"context"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
	"github.com/xuxusheng/time-frequency-be/internal/utils"
)

type IUser interface {
	Create(ctx context.Context, user *model.User) error

	Get(ctx context.Context, id int) (*model.User, error)
	GetByName(ctx context.Context, name string) (*model.User, error)
	ListAndCount(ctx context.Context, p *model.Page, query, role string) ([]*model.User, int, error)
	IsNameExist(ctx context.Context, name string, excludeId int) (bool, error)   // 查询用户名是否被占用
	IsPhoneExist(ctx context.Context, phone string, excludeId int) (bool, error) // 查询手机号是否被占用
	IsEmailExist(ctx context.Context, email string, excludeId int) (bool, error) // 查询邮箱是否被占用

	Update(ctx context.Context, user *model.User, columns []string) error
	UpdatePassword(ctx context.Context, id int, oldPassword, newPassword string) (*model.User, error)

	Delete(ctx context.Context, id int) error
}

func NewUser(dao dao.IUser) *User {
	return &User{Dao: dao}
}

type User struct {
	Dao dao.IUser
}

func (u *User) Create(ctx context.Context, user *model.User) error {
	d := u.Dao
	// 判断用户名是否已存在
	is, err := d.IsNameExist(ctx, user.Name, 0)
	if err != nil {
		return err
	}
	if is {
		return cerror.BadRequest.WithMsg("用户名已存在")
	}

	// 判断手机号是否已存在
	is, err = d.IsPhoneExist(ctx, user.Phone, 0)
	if err != nil {
		return err
	}
	if is {
		return cerror.BadRequest.WithMsg("手机号已存在")
	}

	// 判断邮箱是否已存在
	is, err = d.IsEmailExist(ctx, user.Email, 0)
	if err != nil {
		return err
	}
	if is {
		return cerror.BadRequest.WithMsg("邮箱已存在")
	}

	// 计算密码 hash
	hash, err := utils.EncodePwd(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash

	return d.Create(ctx, user)
}

func (u *User) Get(ctx context.Context, id int) (*model.User, error) {
	return u.Dao.Get(ctx, id)
}

func (u *User) GetByName(ctx context.Context, name string) (*model.User, error) {
	return u.Dao.GetByName(ctx, name)
}

func (u *User) ListAndCount(ctx context.Context, p *model.Page, query, role string) ([]*model.User, int, error) {
	return u.Dao.ListAndCount(ctx, p, query, role)
}

func (u *User) UpdatePhoneAndEmail(ctx context.Context, id int, phone, email string) (*model.User, error) {
	d := u.Dao
	is, err := d.IsPhoneExist(ctx, phone, id)
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

	user := model.User{
		Id:    id,
		Phone: phone,
		Email: email,
	}
	err = d.Update(ctx, &user, []string{"phone", "email"})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) Update(ctx context.Context, user *model.User, columns []string) error {
	// 判断 user.Name 是否已被占用
	if user.Name != "" {
		is, err := u.Dao.IsNameExist(ctx, user.Name, user.Id)
		if err != nil {
			return err
		}
		if is {
			return cerror.BadRequest.WithMsg("用户名已被占用")
		}
	}

	// Phone 是否已被占用
	if user.Phone != "" {
		is, err := u.Dao.IsPhoneExist(ctx, user.Phone, user.Id)
		if err != nil {
			return err
		}
		if is {
			return cerror.BadRequest.WithMsg("手机号已被占用")
		}
	}

	// Email 是否已被占用
	if user.Email != "" {
		is, err := u.Dao.IsEmailExist(ctx, user.Email, user.Id)
		if err != nil {
			return err
		}
		if is {
			return cerror.BadRequest.WithMsg("邮箱已被占用")
		}
	}

	// 计算密码 Hash 值
	if user.Password != "" {
		// 计算新密码 Hash
		hash, err := utils.EncodePwd(user.Password)
		if err != nil {
			return err
		}
		user.Password = hash
	}
	return u.Dao.Update(ctx, user, columns)
}

func (u *User) UpdatePassword(ctx context.Context, id int, oldPassword, newPassword string) (*model.User, error) {
	d := u.Dao

	// 验证旧密码是否正确
	user, err := d.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = utils.ComparePwd(user.Password, oldPassword)
	if err != nil {
		return nil, cerror.BadRequest.WithMsg("旧密码错误")
	}

	// 计算新密码 Hash
	hash, err := utils.EncodePwd(newPassword)
	if err != nil {
		return nil, err
	}

	// 更新密码
	user.Password = hash
	err = d.Update(ctx, user, []string{"password"})
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
