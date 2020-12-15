package dao

import (
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"time"
)

type UserDao struct {
	Dao
}

func NewUserDao(engine *pg.DB) *UserDao {
	return &UserDao{
		Dao{
			engine: engine,
		},
	}
}

func (u *UserDao) Create(name, phone, password string) (*model.User, error) {
	user := model.User{
		Name:     name,
		Phone:    phone,
		Password: password,
	}
	_, err := u.engine.Model(&user).Insert()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserDao) Delete(id uint) error {
	_, err := u.engine.
		Model(&model.User{ID: id}).
		WherePK().
		Delete()
	return err
}

func (u *UserDao) Update(id uint, name, phone string) (*model.User, error) {
	user := model.User{
		ID:        id,
		Name:      name,
		Phone:     phone,
		UpdatedAt: time.Now(),
	}
	// 只更新非零值
	_, err := u.engine.Model(&user).WherePK().Returning("*").UpdateNotZero()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserDao) Get(id uint) (*model.User, error) {
	user := model.User{ID: id}
	err := u.engine.Model(&user).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserDao) ListAndCount(name, phone string, page *model.Page) ([]*model.User, int, error) {
	var users []*model.User
	count, err := u.engine.
		Model(&users).
		Offset(page.Offset()).
		Limit(page.Limit()).
		Where("name LIKE ?", "%"+name+"%").
		Where("phone LIKE ?", "%"+phone+"%").
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (u *UserDao) IsNameExist(name string, excludeID uint) (bool, error) {
	db := u.engine.Model(&model.User{})
	if excludeID != 0 {
		db = db.Where("id != ?", excludeID)
	}
	return db.Where("name = ?", name).Exists()
}

func (u *UserDao) IsPhoneExist(phone string, excludeID uint) (bool, error) {
	db := u.engine.Model(&model.User{})
	if excludeID != 0 {
		db = db.Where("id != ?", excludeID)
	}
	return db.Where("phone = ?", phone).Exists()
}

func (u *UserDao) IsIDExist(id uint) (bool, error) {
	return u.engine.Model(&model.User{ID: id}).WherePK().Exists()
}
