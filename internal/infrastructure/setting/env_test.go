package setting

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/tj/assert"
)

type OrderInfo struct {
	A1 int64   `env:"A1"`
	A2 string  `env:"A2"`
	A3 int     `env:"A3"`
	A4 float32 `env:"A4"`
	A5 float64 `env:"A5"`
	A6 uint    `env:"A6"`
	A7 uint64  `env:"A7"`
	A8 bool    `env:"A8"`
	A9 string  `env:"A9"`
} //env支持的数据类型

func SetEnv() {
	envs := make(map[string]string)
	envs["A1"] = "-1000000000000000000"
	envs["A2"] = "vivi"
	envs["A3"] = "-1011"
	envs["A4"] = "1024.013"
	envs["A5"] = "1024.0000000013"
	envs["A6"] = "1234"
	envs["A7"] = "12340000000000000000"
	envs["A8"] = "false"

	for key := range envs {
		os.Setenv(key, envs[key])
	}
}

func check(act OrderInfo) error {
	data := OrderInfo{
		A1: -1000000000000000000,
		A2: "vivi",
		A3: -1011,
		A4: 1024.013,
		A5: 1024.0000000013,
		A6: 1234,
		A7: 12340000000000000000,
		A8: false,
		A9: "default",
	}
	if data != act {
		return errors.New("actual value in not right")
	}
	return nil
}

//验证
//1.是否能覆盖原来的值，
//2.检测值是否正确，
//3.如果环境变量没有取到，不会覆盖原来的值
func TestFill(t *testing.T) {

	SetEnv()
	order := OrderInfo{
		A1: 0,
		A2: "vivitest",
		A3: -1,
		A4: 1,
		A5: 1,
		A6: 1,
		A7: 1,
		A8: false,
		A9: "default",
	}

	err := FillEnv(&order)
	assert.Nil(t, err)
	fmt.Println(order)
	assert.Nil(t, check(order))
}
