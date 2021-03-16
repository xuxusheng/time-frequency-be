package dao

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/database"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/testdb"
	"log"
	"os"
	"testing"
)

var db *pg.DB

func TestMain(m *testing.M) {

	setup()

	code := m.Run()

	teardown(db)

	os.Exit(code)
}

func setup() {
	var err error

	db, err = database.New(context.Background(), "postgres://postgres:1234@localhost:5432/example2?sslmode=disable")
	if err != nil {
		log.Fatalf("数据库初始化失败：%v", err)
	}

	//db.AddQueryHook(pgdebug.DebugHook{
	//	Verbose: true,
	//})

	// 清空一遍数据，避免被上一次测试异常退出后残留的数据影响
	err = testdb.Truncate(db)
	if err != nil {
		log.Fatalf("清空数据库记录失败：%v", err)
	}
}

// 清空数据库
func truncateTable(db *pg.DB) {
	_ = testdb.Truncate(db)
}

// 打扫战场
func teardown(db *pg.DB) {
	err := testdb.Drop(db)
	if err != nil {
		log.Fatalf("删除数据表失败：%v", err)
	}
	_ = db.Close()
}
