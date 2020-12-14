package dao

import (
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"gorm.io/gorm"
)

type UserDao struct {
	Dao
}

func NewUserDao(engine *gorm.DB) *UserDao {
	// todo 如果 dao 层中，一个 dao 只允许操作单表，那么在 new 的时候，要不要直接把 engine = engine.Model(&model.User{}) 这样固化住？
	return &UserDao{
		Dao{
			engine: engine,
		},
	}
}

// 创建用户
// password 字段需要在 service 层中处理好，转换成 hash
func (u *UserDao) Create(name, phone, password string) (*model.User, error) {
	user := model.User{
		Name:     name,
		Phone:    phone,
		Password: password,
	}
	err := u.engine.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 删除用户，硬删除
func (u *UserDao) Delete(id uint) error {
	// 硬删除
	return u.engine.Unscoped().Delete(&model.User{}, id).Error
}

// 更新用户信息，注意零值也会被更新
func (u *UserDao) Update(id uint, data interface{}) error {
	return u.engine.Model(&model.User{}).Where("id = ?", id).Updates(data).Error
}

// 通过 id 查询单个用户
func (u *UserDao) Get(id uint) (*model.User, error) {
	var user model.User
	err := u.engine.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

// 通过 name、phone 模糊查询命中的用户数量，配合 List 函数实现分页功能
func (u *UserDao) Count(name, phone string) (int64, error) {
	var count int64
	err := u.engine.
		Model(&model.User{}).
		Where("name LIKE ?", "%"+name+"%").
		Where("phone LIKE ?", "%"+phone+"%").
		Count(&count).
		Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 通过 name、phone 模糊查询多个用户列表
func (u *UserDao) List(name, phone string, page *model.Page) ([]*model.User, error) {
	db := u.engine
	var users []*model.User

	// 设置偏移量
	offset := page.Offset()
	if offset >= 0 && page.Ps > 0 {
		db = db.Offset(offset).Limit(page.Ps)
	}

	// 查询
	err := db.
		Where("name LIKE ?", "%"+name+"%").
		Where("phone LIKE ?", "%"+phone+"%").
		Find(&users).
		Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 用户名是否被占用
func (u *UserDao) IsNameExist(name string, excludeID uint) (bool, error) {
	var count int64
	db := u.engine.Model(&model.User{})

	if excludeID != 0 {
		db = db.Where("id != ?", excludeID)
	}

	err := db.Where(model.User{Name: name}).Count(&count).Error
	if err != nil {
		return false, err
	}
	// 数量不为 0 说明存在
	return count != 0, nil
}

// 手机号是否被占用
func (u *UserDao) IsPhoneExist(phone string, excludeID uint) (bool, error) {
	var count int64
	db := u.engine.Model(&model.User{})

	if excludeID != 0 {
		db = db.Where("id != ?", excludeID)
	}

	err := db.Where(model.User{Phone: phone}).Count(&count).Error
	if err != nil {
		return false, err
	}
	// 数量不为 0 说明存在
	return count != 0, nil
}

func (u *UserDao) IsIDExist(id uint) (bool, error) {
	var count int64
	err := u.engine.
		Model(&model.User{}).
		Where(model.User{Model: model.Model{ID: id}}).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}
