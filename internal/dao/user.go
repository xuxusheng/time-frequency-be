package dao

import (
	"context"
	"github.com/go-pg/pg/v10/orm"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"time"
)

type IUser interface {
	// 创建用户
	Create(ctx context.Context, createdBy int, name, nickName, phone, email, password string) (*model.User, error)
	// 获取单个用户
	Get(ctx context.Context, id int) (*model.User, error)
	// 通过 Name 获取用户
	GetByName(ctx context.Context, name string) (*model.User, error)
	// 获取多个用户
	ListAndCount(ctx context.Context, query string, p *model.Page) ([]*model.User, int, error)
	// 更新用户信息
	Update(ctx context.Context, user *model.User, columns []string) error
	// 删除用户
	Delete(ctx context.Context, id int) error
	// 用户名是否存在
	IsNameExist(ctx context.Context, name string, excludeId int) (bool, error)
	// 手机号是否存在
	IsPhoneExist(ctx context.Context, phone string, excludeId int) (bool, error)
	// 邮箱是否存在
	IsEmailExist(ctx context.Context, email string, excludeId int) (bool, error)
}

func NewUser(db orm.DB) *User {
	return &User{
		db: db,
	}
}

type User struct {
	db orm.DB
}

func (u *User) Create(ctx context.Context, createdBy int, name, nickName, phone, email, password string) (*model.User, error) {
	user := model.User{
		Name:        name,
		NickName:    nickName,
		Phone:       phone,
		Email:       email,
		Password:    password,
		CreatedById: createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := u.db.ModelContext(ctx, &user).Returning("*").Insert()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) Get(ctx context.Context, id int) (*model.User, error) {
	user := model.User{Id: id}
	err := u.db.ModelContext(ctx, &user).WherePK().Relation("CreatedBy").Select()
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (u *User) GetByName(ctx context.Context, name string) (*model.User, error) {
	user := model.User{}
	err := u.db.ModelContext(ctx, &user).Where("name = ?", name).Select()
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (u *User) ListAndCount(ctx context.Context, query string, p *model.Page) ([]*model.User, int, error) {
	users := []*model.User{}
	count, err := u.db.ModelContext(ctx, &users).
		Offset(p.Offset()).
		Limit(p.Limit()).
		Where("name LIKE ?", "%"+query+"%").
		WhereOr("phone LIKE ?", "%"+query+"%").
		WhereOr("email LIKE ?", "%"+query+"%").
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return users, count, err
}

func (u *User) Update(ctx context.Context, user *model.User, columns []string) error {
	_, err := u.db.
		ModelContext(ctx, user).
		Column(columns...).
		WherePK().
		Returning("*").
		Update()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(ctx context.Context, id int) error {
	_, err := u.db.ModelContext(ctx, &model.User{Id: id}).WherePK().Delete()
	return err
}

func (u *User) IsNameExist(ctx context.Context, name string, excludeId int) (bool, error) {
	db := u.db.ModelContext(ctx, &model.User{})
	if excludeId != 0 {
		db = db.Where("id != ?", excludeId)
	}
	return db.Where("name = ?", name).Exists()
}

func (u *User) IsPhoneExist(ctx context.Context, phone string, excludeId int) (bool, error) {
	db := u.db.ModelContext(ctx, &model.User{})
	if excludeId != 0 {
		db = db.Where("id != ?", excludeId)
	}
	return db.Where("phone = ?", phone).Exists()
}

func (u *User) IsEmailExist(ctx context.Context, email string, excludeId int) (bool, error) {
	db := u.db.ModelContext(ctx, &model.User{})
	if excludeId != 0 {
		db = db.Where("id != ?", excludeId)
	}
	return db.Where("email = ?", email).Exists()
}
