package model

import "time"

type User struct {
	ID       uint   `json:"id" example:"1"`                                        // 用户 ID
	Name     string `json:"name" example:"xusheng" pg:",unique,notnull"`           // 用户名
	Phone    string `json:"phone" example:"17707272442" pg:",unique,notnull"`      // 手机号
	Role     string `json:"role" example:"member" pg:"default:('member'),notnull"` // 用户角色
	Password string `json:"-" pg:",notnull"`                                       // 用户密码

	CreatedAt time.Time `json:"created_at" example:"2020-12-09T18:52:41.555+08:00" pg:"default:now(),notnull"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" example:"2020-12-09T18:52:41.555+08:00" pg:"default:now(),notnull"` // 更新时间
}
