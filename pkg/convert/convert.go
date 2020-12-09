package convert

import "strconv"

// 将 string 类型转为其他类型时使用
type StrTo string

func (s StrTo) String() string {
	return string(s)
}

// string 转换为 int 类型
func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

// 一定可以转化为 int 类型，忽略掉错误
func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) UInt() (uint, error) {
	v, err := s.Int()
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func (s StrTo) MustUInt() uint {
	v, _ := s.UInt()
	return v
}

func (s StrTo) UInt32() (uint32, error) {
	v, err := s.Int()
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

func (s StrTo) MustInt32() uint32 {
	v, _ := s.UInt32()
	return v
}
