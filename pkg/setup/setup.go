package setup

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
func Logger() {
	global.Logger = logger.NewLogger(os.Stdout, "", log.LstdFlags)
}

// 准备全局的配置
func Setting() error {

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
	err = s.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("PG", &global.PGSetting)
	if err != nil {
		return err
	}

	// 从环境变量中读取一部分配置，优先级大于配置文件，小于启动命令参数
	// todo 这里可以看看 viper 有没有提供什么简单的从环境变量覆盖配置文件的功能，然后优化一下

	if mode := os.Getenv("SERVER_MODE"); mode != "" {
		global.ServerSetting.Mode = mode
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		global.ServerSetting.HttpPort = port
	}
	if pgDBName := os.Getenv("PG_DBNAME"); pgDBName != "" {
		global.PGSetting.DBName = pgDBName
	}
	if pgUsername := os.Getenv("PG_USERNAME"); pgUsername != "" {
		global.PGSetting.Username = pgUsername
	}
	if pgPassword := os.Getenv("PG_PASSWORD"); pgPassword != "" {
		global.PGSetting.Password = pgPassword
	}
	if pgHost := os.Getenv("PG_HOST"); pgHost != "" {
		global.PGSetting.Host = pgHost
	}

	// 有这么个从启动命令参数中取配置的功能，但是以目前自己常用的部署方案来说，没啥必要支持这个
	// 从启动命令参数中取
	//flag.StringVar(&port, "port", "", "服务器监听端口")
	//flag.StringVar(&runMode, "mode", "", "启动模式，debug 或 release")
	//if port != "" {
	//	global.ServerSetting.HttpPort = port
	//}

	// 默认从 yaml 文件中导入进来的时间，单位不是秒，需要转换一下
	global.JWTSetting.Expire *= time.Second
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	return nil
}

func PGEngine(pgSetting *setting.PGSettingS, serverSetting *setting.ServerSettingS) error {
	var err error
	global.PGEngine, err = model.NewPGEngine(pgSetting, serverSetting.Mode)
	return err
}

func Reset() {
	global.Logger = nil
	global.ServerSetting = nil
	global.PGSetting = nil
	global.AppSetting = nil
	global.JWTSetting = nil
	global.PGEngine = nil
}
