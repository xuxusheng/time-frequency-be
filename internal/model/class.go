package model

import "time"

// 班级表
type Class struct {
	// --- 表名 ---
	tableName struct{} `pg:"class"`

	// --- 业务字段 ---
	Name        string `json:"name" pg:",unique,notnull"`                     // 班级名称
	Description string `json:"description" pg:",use_zero,notnull,default:''"` // 班级描述

	// --- 关联字段 ---
	Members     []*User `json:"members" pg:"rel:has-many"`   // 班级包含的成员
	CreatedById int     `json:"-" pg:",notnull"`             // 创建人ID
	CreatedBy   *User   `json:"created_by" pg:"rel:has-one"` // 创建人

	// --- 通用字段 ---
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at" pg:",notnull,default:now()"`
	UpdatedAt time.Time `json:"updated_at" pg:",notnull,default:now()"`
}
