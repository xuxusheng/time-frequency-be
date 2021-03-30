package model

import "time"

// 用户表
type User struct {
	// --- 表名 ---
	tableName struct{} `pg:"user"`

	// --- 业务字段 ---
	Name     string `json:"name" pg:",unique,notnull"`  // 用户名
	Phone    string `json:"phone" pg:",unique,notnull"` // 手机号
	Email    string `json:"email" pg:",unique,notnull"` // 邮箱
	Password string `json:"-" pg:",notnull"`

	// --- 关联字段 ---
	ClassId int    `json:"-"`
	Class   *Class `json:"class" pg:"rel:has-one"` // 用户所属的班级

	CreatedById int   `json:"-"`
	CreatedBy   *User `json:"created_by" pg:"rel:has-one"` // 用户的创建人

	LearningMaterials []*LearningMaterial `json:"-" pg:"rel:has-many"` // 由用户上传的学习资料

	// --- 通用字段 ---
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at" pg:",notnull,default:now()"`
	UpdatedAt time.Time `json:"updated_at" pg:",notnull,default:now()"`
}
