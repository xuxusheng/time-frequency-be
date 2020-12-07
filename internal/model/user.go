package model

import (
	"github.com/xuxusheng/time-frequency-be/global"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null;unique"`         // 用户名（唯一、非空）
	Phone    string `json:"phone" gorm:"not null;unique"`        // 手机号（唯一、非空）
	Role     string `json:"role" gorm:"not null;default:member"` // 用户角色，admin & member（非空、默认 member）
	Password string `json:"-" gorm:"not null"`                   // 用户密码，计算 hash 后存入（非空）
}

func (u User) TableName() string {
	return global.DatabaseSetting.TablePrefix + "user"
}

// 计算匹配到的 user 数量，通过 name 和 phone 模糊匹配
func (u User) Count(db *gorm.DB) (int64, error) {
	var count int64
	err := db.
		Model(&u).
		Where("name LIKE ?", "%"+u.Name+"%").
		Where("phone LIKE ?", "%"+u.Phone+"%").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 批量查询 user，通过 name 和 phone 模糊匹配
func (u User) List(db *gorm.DB, pageOffset, pageSize int) ([]*User, error) {
	var users []*User
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}

	err := db.
		Where("name LIKE ?", "%"+u.Name+"%").
		Where("phone LIKE ?", "%"+u.Phone+"%").
		Find(&users).
		Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 新增用户
func (u User) Create(db *gorm.DB) error {
	return db.Create(&u).Error
}

// 更新用户
func (u User) Update(db *gorm.DB, values interface{}) error {
	return db.Model(&u).Where("id = ?", u.ID).Updates(values).Error
}

// 删除用户
func (u User) Delete(db *gorm.DB) error {
	// 硬删除
	return db.Unscoped().Delete(&u).Error
}

// 获取单个用户
func (u User) Get(db *gorm.DB) (*User, error) {
	//var user User
	//if u.ID != 0 {
	//	db = db.Where("id = ?", u.ID)
	//}
	//if u.Name != "" {
	//	db = db.Where("name = ?", u.Name)
	//}
	//if u.Phone != "" {
	//	db = db.Where("phone = ?", u.Phone)
	//}
	//err := db.First(&user).Error
	//return &user, err

	var user User
	// 此种方式只会针对 u 结构体中非零值字段进行查询
	err := db.Where(&u).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
