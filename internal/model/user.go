package model

import (
	"github.com/xuxusheng/time-frequency-be/global"
)

type User struct {
	Model
	Name     string `json:"name" gorm:"not null;unique"`         // 用户名（唯一、非空）
	Phone    string `json:"phone" gorm:"not null;unique"`        // 手机号（唯一、非空）
	Role     string `json:"role" gorm:"not null;default:member"` // 用户角色，admin & member（非空、默认 member）
	Password string `json:"-" gorm:"not null"`                   // 用户密码，计算 hash 后存入（非空）
}

func (u User) TableName() string {
	return global.DatabaseSetting.TablePrefix + "user"
}
