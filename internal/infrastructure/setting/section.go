package setting

import "time"

type Server struct {
	Mode         string        // 服务器运行模式 debug、test 或 release
	HttpPort     int           // 服务器监听端口
	ReadTimeout  time.Duration // request 超时时间
	WriteTimeout time.Duration // response 超时时间
}

type App struct {
	DefaultPs int // 默认每页查询记录条数
	MaxPs     int // 每页最多查询记录条数
}

type JWT struct {
	Secret string        // JWT 密钥
	Issuer string        // 签发人
	Expire time.Duration // Token 过期时间
}

type DB struct {
	Host     string // 数据库地址
	Database string // 数据库名称
	User     string // 用户名
	Password string // 密码
}
