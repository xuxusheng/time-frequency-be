package service

import (
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
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

	err := testdb.Drop(db)
	if err != nil {
		log.Fatalf("测试完成后删除数据表失败：%v", err)
	}
	os.Exit(code)
}

func setup() {
	var err error
	db, err = testdb.New()
	if err != nil {
		log.Fatalf("测试数据库连接失败：%v", err)
	}

	userDao = dao.NewUser(db)
	subjectDao = dao.NewSubject(db)
	classDao = dao.NewClass(db)
	lmDao = dao.NewLearningMaterial(db)
}
