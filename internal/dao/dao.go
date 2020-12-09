package dao

import "gorm.io/gorm"

type Dao struct {
	engine *gorm.DB
	// 如果需要的话，可以把 gin.Context 也挂进来
}
