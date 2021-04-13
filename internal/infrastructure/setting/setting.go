package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type Setting struct {
	vp *viper.Viper

	Server *Server
	App    *App
	JWT    *JWT
	DB     *DB
}

func New(configPath ...string) (*Setting, error) {
	vp := viper.New()

	// 设置配置文件名称为 config
	vp.SetConfigName("config")
	vp.AddConfigPath("config/")

	for _, path := range configPath {
		if path != "" {
			vp.AddConfigPath(path)
		}
	}

	// 设置配置文件类型为 yaml
	vp.SetConfigType("yaml")

	// 读取配置文件
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	s := &Setting{vp: vp}

	// 监听配置文件变化
	s.Watch()

	return s, nil
}

func (s *Setting) Init() error {
	vp := s.vp

	err := vp.UnmarshalKey("Server", &s.Server)
	if err != nil {
		return err
	}

	err = vp.UnmarshalKey("App", &s.App)
	if err != nil {
		return err
	}

	err = vp.UnmarshalKey("JWT", &s.JWT)
	if err != nil {
		return err
	}

	err = vp.UnmarshalKey("DB", &s.DB)
	if err != nil {
		return err
	}

	// 读取系统环境变量
	err = FillEnv(s.Server)
	if err != nil {
		return err
	}
	err = FillEnv(s.App)
	if err != nil {
		return err
	}
	err = FillEnv(s.DB)
	if err != nil {
		return err
	}
	err = FillEnv(s.JWT)
	if err != nil {
		return err
	}

	return nil
}

func (s *Setting) Watch() {
	go func() {
		s.vp.WatchConfig()
		// 重新初始化
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.Init()
			log.Println("重新加载配置项 Ready！👌")
		})
	}()
}
