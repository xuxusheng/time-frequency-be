package main

import (
	"context"
	"flag"
	"github.com/kataras/iris/v12"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/router"
	"github.com/xuxusheng/time-frequency-be/pkg/logger"
	"github.com/xuxusheng/time-frequency-be/pkg/setting"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// @title 时频学习平台
// @version 1.0
// @description 时频学习平台后端接口文档

// @contact.name xusheng
// @contact.url https://github.com/xuxusheng
// @contact.email 20691718@qq.com
func main() {
	app := router.NewApp()

	go func() {

		log.Println("发射！🚀")

		s := &http.Server{
			Addr:           ":" + global.ServerSetting.HttpPort,
			Handler:        app,
			ReadTimeout:    global.ServerSetting.ReadTimeout,
			WriteTimeout:   global.ServerSetting.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		}

		err := app.Run(
			iris.Server(s),
			iris.WithOptimizations, // 开启优化功能，比如压缩返回的 json 字符串之类的
			iris.WithoutServerError(iris.ErrServerClosed), // 忽略掉服务器关闭错误
		)

		if err != nil {
			log.Fatalf("发射失败 ☠️ : %v", err)
		}
	}()

	// 等待终端信息
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("返航中...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatal("返航失败，强制着陆 🙏 : ", err)
	}
	log.Println("返航成功，拜拜~ 👋")

}

func init() {
	setupLogger()
	log.Println("日志组件 Ready! 👌")

	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	log.Println("配置项 Ready! 👌")

	err = setupPGEngine()
	if err != nil {
		log.Fatalf("init.setupPGEngine err: %v", err)
	}
	log.Println("数据库连接 Ready! 👌")
}

// 准备 Logger
func setupLogger() {
	global.Logger = logger.NewLogger(os.Stdout, "", log.LstdFlags)
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

func setupPGEngine() error {
	var err error
	global.PGEngine, err = model.NewPGEngine(global.PGSetting, global.ServerSetting.Mode)
	return err
}
