package v1

import (
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/testdb"
	"github.com/xuxusheng/time-frequency-be/internal/service"
	"log"
	"os"
	"testing"
)

var db *pg.DB
var userSvc *service.User

func TestMain(m *testing.M) {
	setup()

	userSvc = service.NewUser(dao.NewUser(db))

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
}
