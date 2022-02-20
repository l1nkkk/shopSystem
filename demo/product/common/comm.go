package common

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

// l1nkkk: 充分体现了代码复用的特点，不然要给每一个存入到数据库的对象，定义一个转化，
// 这样的话很灵活，但是代码复杂

// 过程： obj 和 data 之间的桥梁为 tag:`sql`,遍历 obj中的所有属性，
// 并找到每个属性在 data 中的对应的值 val_str（string类型）,
// 找到每个属性对应的类型名 typename_str, 通过(val_str, typename_str)
// 将val_str转化成对应的reflect.value


// DataToStructByTagSql 根据结构体中sql标签映射数据到结构体中并且转换类型
func DataToStructByTagSql(data map[string]string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()

	// 获取属性个数 objValue.NumField()
	for i := 0; i < objValue.NumField(); i++ {
		// 第一次循环获取map值为： data["ID"]
		// 获取sql对应的值，该值是string了下
		value := data[objValue.Type().Field(i).Tag.Get("sql")]
		// 获取对应字段的名称，如ID
		name := objValue.Type().Field(i).Name
		// 获取对应字段类型，如int64
		structFieldType := objValue.Field(i).Type()
		// 获取变量类型，也可以直接写"string类型"
		val := reflect.ValueOf(value)
		var err error
		if structFieldType != val.Type() {
			// 类型转换，对string类型的值进行转化，转化成原始的类型
			val, err = TypeConversion(value, structFieldType.Name()) //类型转换
			if err != nil {

			}
		}
		//设置类型值
		objValue.FieldByName(name).Set(val)
	}
}

// TypeConversion类型转换
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

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}

