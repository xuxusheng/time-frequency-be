package dao

import (
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/pkg/app"
	"gorm.io/gorm"
)

func (d *Dao) CountUser(name, phone string) (int64, error) {
	user := model.User{
		Name:  name,
		Phone: phone,
	}
	return user.Count(d.engine)
}

func (d *Dao) GetUserList(name, phone string, pn, ps int) ([]*model.User, error) {
	user := model.User{
		Name:  name,
		Phone: phone,
	}
	offset := app.GetPageOffset(pn, ps)
	return user.List(d.engine, offset, ps)
}

func (d *Dao) CreateUser(name, phone, password string) error {
	user := model.User{
		Name:     name,
		Phone:    phone,
		Password: password,
	}
	return user.Create(d.engine)
}

func (d *Dao) GetUser(id uint, name, phone string) (*model.User, error) {
	user := model.User{
		Model: gorm.Model{
			ID: id,
		},
		Name:  name,
		Phone: phone,
	}

	err := user.Get(d.engine)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
