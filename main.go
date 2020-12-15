package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
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
	gin.SetMode(global.ServerSetting.RunMode)
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router.NewRouter(),
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err: %v", err)
		}
	}()

	// 等待终端信息
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")

}

func init() {
	setupLogger()

	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupPGEngine()
	if err != nil {
		log.Fatalf("init.setupPGEngine err: %v", err)
	}
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

	var runMode string

	// 从环境变量中读取一部分配置，优先级大于配置文件，小于启动命令参数
	// todo 这里可以看看 viper 有没有提供什么简单的从环境变量覆盖配置文件的功能，然后优化一下

	if port := os.Getenv("SERVER_PORT"); port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode = os.Getenv("SERVER_MODE"); runMode != "" {
		global.ServerSetting.RunMode = runMode
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
	//if runMode != "" {
	//	global.ServerSetting.RunMode = runMode
	//}

	// 默认从 yaml 文件中导入进来的时间，单位不是秒，需要转换一下
	global.JWTSetting.Expire *= time.Second
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	return nil

}

func setupPGEngine() error {
	var err error
	global.PGEngine, err = model.NewPGEngine(global.PGSetting)
	return err
}
