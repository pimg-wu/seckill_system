package common

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

//根据结构体中的sql标签映射数据到结构体中并转换类型
func DataToStructByTagSql(data map[string]string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		//获取sql对应的值
		value := data[objValue.Type().Field(i).Tag.Get("sql")]
		//获取对应字段名称
		name := objValue.Type().Field(i).Name
		//获取对应字段类型
		structFieldType := objValue.Field(i).Type()
		//获取变量类型
		val := reflect.ValueOf(value)
		var err error
		if structFieldType != val.Type() {
			//类型转换
			val, err = TypeConversion(value, structFieldType.Name())
			if err != nil {

			}
		}
		//设置类型值
		objValue.FieldByName(name).Set(val)
	}
}

func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	return reflect.ValueOf(value), errors.New("未知类型：" + ntype)
}
