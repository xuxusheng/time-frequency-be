package v1

import (
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/testdb"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"testing"
)

// 准备用户数据
func prepareUser(t *testing.T, db *pg.DB) []*model.User {
	users, err := testdb.SeedUser(db)
	if err != nil && len(users) == 0 {
		t.Fatalf("准备用户数据失败：%v", err)
	}
	return users
}

func TestUser_Create(t *testing.T) {
	// 准备数据

	// 初始化服务

	// 绑定路由

	// 发送请求

	// 验证结果
}
