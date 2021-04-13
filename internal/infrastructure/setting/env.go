package setting

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func upper(v string) string {
	return strings.ToUpper(v)
}

func fill(ind reflect.Value) error {
	for i := 0; i < ind.NumField(); i++ {
		f := ind.Type().Field(i)
		name := f.Name
		envName, exist := f.Tag.Lookup("env")
		if exist {
			name = envName
		} else {
			name = upper(name)
		}
		err := parse(name, ind.Field(i), f)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseBool(v string) (bool, error) {
	if v == "" {
		return false, nil
	}
	return strconv.ParseBool(v)
}

func parse(prefix string, f reflect.Value, sf reflect.StructField) error {
	//df := sf.Tag.Get("default")
	/*	isRequire, err := parseBool(sf.Tag.Get("require"))
		if err != nil {
			return fmt.Errorf("the value of %s is not a valid `member` of bool ，only "+
				"[1 0 t f T F true false TRUE FALSE True False] are supported", prefix)
		}*/
	ev, exist := os.LookupEnv(prefix)

	/*if !exist && isRequire {
		return fmt.Errorf("%s is required, but has not been set", prefix)
	}*/
	/*if !exist && df != "" {
		ev = df
	}*/

	if !exist { //如果不存在则跳过
		return nil
	}
	//log.Print("ev:", ev)
	switch f.Kind() {
	case reflect.String:
		f.SetString(ev)
	case reflect.Int:
		iv, err := strconv.ParseInt(ev, 10, 32)
		if err != nil {
			return fmt.Errorf("%s:%w", prefix, err)
		}
		f.SetInt(iv)
	case reflect.Int64:
		if f.Type().String() == "time.Duration" {
			t, err := time.ParseDuration(ev)
			if err != nil {
				return fmt.Errorf("%s:%w", prefix, err)
			}
			f.Set(reflect.ValueOf(t))
		} else {
			iv, err := strconv.ParseInt(ev, 10, 64)
			if err != nil {
				return fmt.Errorf("%s:%w", prefix, err)
			}
			f.SetInt(iv)
		}
	case reflect.Uint:
		uiv, err := strconv.ParseUint(ev, 10, 32)
		if err != nil {
			return fmt.Errorf("%s:%w", prefix, err)
		}
		f.SetUint(uiv)
	case reflect.Uint64:
		uiv, err := strconv.ParseUint(ev, 10, 64)
		if err != nil {
			return fmt.Errorf("%s:%w", prefix, err)
		}
		f.SetUint(uiv)
	case reflect.Float32:
		f32, err := strconv.ParseFloat(ev, 32)
		if err != nil {
			return fmt.Errorf("%s:%w", prefix, err)
		}
		f.SetFloat(f32)
	case reflect.Float64:
		f64, err := strconv.ParseFloat(ev, 64)
		if err != nil {
			return fmt.Errorf("%s:%w", prefix, err)
		}
		f.SetFloat(f64)
	case reflect.Bool:
		b, err := parseBool(ev)
		if err != nil {
			return fmt.Errorf("%s:%w", prefix, err)
		}
		f.SetBool(b)
	}
	return nil
}

func FillEnv(v interface{}) error {
	ind := reflect.Indirect(reflect.ValueOf(v))
	if reflect.ValueOf(v).Kind() != reflect.Ptr || ind.Kind() != reflect.Struct {
		return fmt.Errorf("only the pointer to a struct is supported")
	}
	err := fill(ind)
	if err != nil {
		return err
	}
	return nil
}
