package global

import (
	"github.com/xuxusheng/time-frequency-be/pkg/setting"
)

var (
	// 全局配置
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	JWTSetting      *setting.JWTSettingS
)
