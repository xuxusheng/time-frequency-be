package setting

import "time"

type ServerSettingS struct {
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type AppSettingS struct {
	DefaultPageSize int
	MaxPageSize     int
}

type JWTSettingS struct {
	Secret string
	Issuer string
	Expire time.Duration
}

// 用来存储已经读取的配置
var sections = make(map[string]interface{})

func (s *Setting) ReadSection(k string, v interface{}) error {
	// 通过 key，从 vp 从查到相应的字符串，然后 parse 成对象
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}

func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}

	// todo 重载完之后，需要把几个时间单位重设一下，可以封一个函数
	return nil
}
