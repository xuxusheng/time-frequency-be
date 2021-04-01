package model

import "time"

// 科目表
type Subject struct {
	// --- 表名 ---
	tableName struct{} `pg:"subject"`

	// --- 业务字段 ---
	Name        string `json:"name" pg:",unique,notnull"`                     // 科目名称
	Description string `json:"description" pg:",use_zero,notnull,default:''"` // 科目描述

	// --- 关联字段
	LearningMaterials []*LearningMaterial `json:"-" pg:"rel:has-many"` // 科目下包含的所有学习资料
	CreatedById       int                 `json:"-" pg:",notnull"`     // 创建人ID
	CreatedBy         *User               `json:"-" pg:"rel:has-one"`  // 创建人

	// --- 通用字段 ---
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at" pg:",notnull,default:now()"`
	UpdatedAt time.Time `json:"updated_at" pg:",notnull,default:now()"`
}
