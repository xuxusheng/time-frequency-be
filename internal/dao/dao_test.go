package dao

import (
	"context"
	"github.com/fatih/color"
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/pkg/setting"
	"github.com/xuxusheng/time-frequency-be/pkg/setup"
	"os"
	"testing"
)

var DB *pg.DB

func TestMain(t *testing.M) {
	color.Cyan("dao_test.testMain 函数开始")

	var err error
	DB, err = model.NewPGEngine(&setting.PGSettingS{
		Host:        "127.0.0.1:5432",
		DBName:      "gotest",
		Username:    "gotest",
		TablePrefix: "gotest",
		Password:    "gotest",
	}, "test")
	if err != nil {
		panic(err)
	}

	code := t.Run()
	tearDown()
	color.Cyan("dao_test.testMain 函数结束")
	os.Exit(code)
}

func tearDown() {
	setup.Reset()
}

func TestDBVersion(t *testing.T) {
	err := DB.Ping(context.Background())
	if err != nil {
		t.Errorf("db ping failed: %v", err)
	}
}
