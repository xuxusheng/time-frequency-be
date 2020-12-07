package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

func NewSetting(configPath ...string) (*Setting, error) {
	vp := viper.New()
	// 设置配置文件名称为 config
	vp.SetConfigName("config")

	for _, path := range configPath {
		if path != "" {
			vp.AddConfigPath(path)
		}
	}

	// 设置配置文件类型为 yaml
	vp.SetConfigType("yaml")

	// 读取配置文件
	// yaml 是一个多级嵌套的结构，这里读取进来是一级的key，然后对应的值是字符串，后面再将字符串 parse 成对象
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	s := &Setting{vp}
	s.WatchSettingChange()
	return s, nil
}

func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSection()
		})
	}()
}
