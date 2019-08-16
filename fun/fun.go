package fun

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/sony/sonyflake"
	"reflect"
	"strconv"
	"strings"
)

var sf = sonyflake.NewSonyflake(sonyflake.Settings{})

func Unique() (string, error) {
	id, err := sf.NextID()
	if err != nil {
		return ``, err
	}
	return strconv.FormatUint(id, 10), nil
}

//传统md5
func Md5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

//转为下划线小写
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

//转为驼峰
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

//value is nil
func IsNil(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

//copy struct or *struct
func Clone(obj interface{}) (i interface{}) {
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		// Pointer:
		i = reflect.New(reflect.ValueOf(obj).Elem().Type()).Interface()
	} else {
		// Not pointer:
		i = reflect.New(reflect.TypeOf(obj)).Elem().Interface()
	}
	return
}

//copy and "make"
func MakeClone(data interface{}) interface{} {

	if reflect.ValueOf(data).Elem().Len() == 0 {
		t := reflect.ValueOf(data).Elem().Type()
		data := reflect.New(reflect.ValueOf(data).Elem().Type())
		data.Elem().Set(reflect.MakeSlice(t, 0, 0))
		return data.Interface()
	}
	return data
}
