package pg_dao

import (
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/pg_model"
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

func (u *UserDao) Create(name, phone, password string) (*pg_model.User, error) {
	user := pg_model.User{
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

func (u *UserDao) Delete(id int) error {
	_, err := u.engine.
		Model(&pg_model.User{ID: id}).
		WherePK().
		Delete()
	return err
}

func (u *UserDao) Update(id int, name, phone string) error {
	user := pg_model.User{
		ID:    id,
		Name:  name,
		Phone: phone,
	}
	// 只更新非零值
	_, err := u.engine.Model(&user).WherePK().UpdateNotZero()
	return err
}

func (u *UserDao) Get(id int) (*pg_model.User, error) {
	user := pg_model.User{ID: id}
	err := u.engine.Model(&user).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserDao) ListAndCount(name, phone string, page *model.Page) ([]*pg_model.User, int, error) {
	var users []pg_model.User
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
	var newUsers []*pg_model.User
	for _, v := range users {
		newUsers = append(newUsers, &v)
	}
	return newUsers, count, nil
}

func (u *UserDao) IsNameExist(name string, excludeID int) (bool, error) {
	db := u.engine.Model(&pg_model.User{})
	if excludeID != 0 {
		db = db.Where("id != ?", excludeID)
	}
	return db.Where("name = ?", name).Exists()
}

func (u *UserDao) IsPhoneExist(phone string, excludeID int) (bool, error) {
	db := u.engine.Model(&pg_model.User{})
	if excludeID != 0 {
		db = db.Where("id != ?", excludeID)
	}
	return db.Where("phone = ?", phone).Exists()
}

func (u *UserDao) IsIDExist(id int) (bool, error) {
	return u.engine.Model(&pg_model.User{ID: id}).WherePK().Exists()
}
