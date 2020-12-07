package main

import (
	"flag"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/pkg/logger"
	"github.com/xuxusheng/time-frequency-be/pkg/setting"
	"log"
	"os"
	"strings"
	"time"
)

// 准备 Logger
func setupLogger() {
	global.Logger = logger.NewLogger(os.Stdout, "", log.LstdFlags).WithCaller(2)
}

// 准备全局的配置
func setupSetting() error {

	var configPath string

	flag.StringVar(&configPath, "configPath", "config/", "配置文件存放路径，多个路径用英文逗号分隔")

	s, err := setting.NewSetting(strings.Split(configPath, ",")...)
	if err != nil {
		return err
	}

	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	var runMode string

	// 从环境变量中读取一部分配置，优先级大于配置文件，小于启动命令参数
	// todo 这里可以看看 viper 有没有提供什么简单的从环境变量覆盖配置文件的功能，然后优化一下

	if port := os.Getenv("SERVER_PORT"); port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode = os.Getenv("SERVER_MODE"); runMode != "" {
		global.ServerSetting.RunMode = runMode
	}
	if DBType := os.Getenv("DB_TYPE"); DBType != "" {
		global.DatabaseSetting.DBType = DBType
	}
	if DBUser := os.Getenv("DB_USER"); DBUser != "" {
		global.DatabaseSetting.UserName = DBUser
	}
	if DBPwd := os.Getenv("DB_PWD"); DBPwd != "" {
		global.DatabaseSetting.Password = DBPwd
	}
	if DBHost := os.Getenv("DB_HOST"); DBHost != "" {
		global.DatabaseSetting.Host = DBHost
	}
	if DBName := os.Getenv("DB_Name"); DBName != "" {
		global.DatabaseSetting.DBName = DBName
	}
	if tablePrefix := os.Getenv("DB_TABLE_PREFIX"); tablePrefix != "" {
		global.DatabaseSetting.TablePrefix = tablePrefix
	}

	// 有这么个从启动命令参数中取配置的功能，但是以目前自己常用的部署方案来说，没啥必要支持这个
	// 从启动命令参数中取
	//flag.StringVar(&port, "port", "", "服务器监听端口")
	//flag.StringVar(&runMode, "mode", "", "启动模式，debug 或 release")
	//if port != "" {
	//	global.ServerSetting.HttpPort = port
	//}
	//if runMode != "" {
	//	global.ServerSetting.RunMode = runMode
	//}

	// 默认从 yaml 文件中导入进来的时间，单位不是秒，需要转换一下
	global.JWTSetting.Expire *= time.Second
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	return nil

}

// 准备数据库连接
func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}
