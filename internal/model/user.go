package model

import "time"

type User struct {
	ID       uint   `json:"id" example:"1"`
	Name     string `json:"name" example:"xusheng" pg:",notnull"`
	Phone    string `json:"phone" example:"17707272442" pg:",notnull"`
	Role     string `json:"role" example:"member" pg:"default:('member'),notnull"`
	Password string `json:"-" pg:",notnull"`

	CreatedAt time.Time `json:"created_at" example:"2020-12-09T18:52:41.555+08:00" pg:"default:now(),notnull"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-12-09T18:52:41.555+08:00" pg:"default:now(),notnull"`
}
