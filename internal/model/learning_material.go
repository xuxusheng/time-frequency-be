package model

import "time"

// 学习资料
type LearningMaterial struct {
	// --- 表名 ---
	tableName struct{} `pg:"learning_material"`

	// --- 业务字段 ---
	Name        string `json:"name" pg:",unique,notnull"`                     // 资料名称
	Description string `json:"description" pg:",use_zero,notnull,default:''"` // 资料描述
	Md5         string `json:"-" pg:",notnull"`
	FilePath    string `json:"-" pg:",notnull"` // 文件存放的路径

	// --- 关联字段 ---
	SubjectId int      `json:"-"`
	Subject   *Subject `json:"-" pg:"rel:has-one"` // 资料所属的科目

	CreatedById int   `json:"-" pg:",notnull"`    // 资料的创建人 ID
	CreatedBy   *User `json:"-" pg:"rel:has-one"` // 资料的创建人

	UpdatedById int   `json:"-" pg:",notnull"`    // 资料的更新人 ID
	UpdatedBy   *User `json:"-" pg:"rel:has-one"` // 资料的更新人

	// --- 通用字段 ---
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at" pg:",notnull,default:now()"`
	UpdatedAt time.Time `json:"updated_at" pg:",notnull,default:now()"`
}
