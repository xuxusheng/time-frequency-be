package service

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/database"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/testdb"
	"log"
	"os"
	"testing"
)

var db *pg.DB
var userDao *dao.User
var subjectDao *dao.Subject
var classDao *dao.Class
var lmDao *dao.LearningMaterial

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
		log.Fatalf("数据库初始化连接失败：%v", err)
	}
	userDao = dao.NewUser(db)
	subjectDao = dao.NewSubject(db)
	classDao = dao.NewClass(db)
	lmDao = dao.NewLearningMaterial(db)

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
