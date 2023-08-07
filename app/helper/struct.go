/*
 * @Description:
 * @Author: mali
 * @Date: 2023-07-31 15:14:58
 * @LastEditTime: 2023-07-31 15:16:06
 * @LastEditors: VSCode
 * @Reference:
 */
package helper

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

/**
 * @Author: mali
 * @Func:
 * @Description: 结构体转化为bson.M
 * @Param:
 * @Return:
 * @Example:
 * @param {context.Context} ctx
 * @param {interface{}} param
 */
func StructToBson(ctx context.Context, param interface{}) (bson.M, error) {
	data, _ := bson.Marshal(param)
	mmap := bson.M{}
	err := bson.Unmarshal(data, mmap)
	if err != nil {
		return nil, err
	}
	return mmap, nil
}

/**
 * len(indexKey) > 0 将切片类型的结构体，转成 map 输出，输出结果存在了 desk 中，也就是你的第一个参数
 * len(indexKey) == 0 将切片类型的结构体，转成 slice 输出，输出结果存在了 desk 中，也就是你的第一个参数
 *
 * @demo 假如有下述结构体，
 * type User struct {
 *  ID   int
 *  NAME string
 * }
 * 输入：
 * @params desk &map[int64]User   desk是个指针呦！！！
 * @params input []User{User{ID:1, NAME:"zwk"}, User{ID:2, NAME:"zzz"}}
 * @params indexKey 键名 "ID" 只支持的索引类型为 (unint unint8 unint16 unint32 unint64 int int8 int16 int32 int64 float32 float64)
 * @params columnKey 列名 空字符串代表返回整个结构体，反之返回结构体中的某一列
 *
 * 输出：
 * err 错误信息
 * 入参 desk 已经被赋值：map[int]User{1:User{ID:1, NAME:"zwk"}, 2:User{ID:2, NAME:"zzz"}}
 *
 * 输入：
 * @params desk &[]int   desk是个指针呦！！！
 * @params input []User{User{ID:1, NAME:"zwk"}, User{ID:2, NAME:"zzz"}}
 * @params indexKey 键名 ""
 * @params columnKey 列名
 *
 * 输出：
 * err 错误信息
 * 入参 desk 已经被赋值：[]int{1, 2}
 */
func StructColumn(desk, input interface{}, columnKey, indexKey string) (err error) {
	structIndexColumn := func(desk, input interface{}, columnKey, indexKey string) (err error) {
		findStructValByIndexKey := func(curVal reflect.Value, elemType reflect.Type, indexKey, columnKey string) (indexVal, columnVal reflect.Value, err error) {
			indexExist := false
			columnExist := false
			for i := 0; i < elemType.NumField(); i++ {
				curField := curVal.Field(i)
				if elemType.Field(i).Name == indexKey {
					switch curField.Kind() {
					case reflect.String, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int, reflect.Float64, reflect.Float32:
						indexExist = true
						indexVal = curField
					case reflect.Struct:

					default:
						return indexVal, columnVal, fmt.Errorf("indexKey must be unint int float or string")
					}
				}
				if elemType.Field(i).Name == columnKey {
					columnExist = true
					columnVal = curField
					continue
				}
			}
			if !indexExist {
				return indexVal, columnVal, fmt.Errorf("indexKey %s not found in %s's field", indexKey, elemType)
			}
			if len(columnKey) > 0 && !columnExist {
				return indexVal, columnVal, fmt.Errorf("columnKey %s not found in %s's field", columnKey, elemType)
			}
			return
		}
		deskValue := reflect.ValueOf(desk)
		if deskValue.Elem().Kind() != reflect.Map {
			return fmt.Errorf("desk must be map")
		}
		deskElem := deskValue.Type().Elem()
		if len(columnKey) == 0 && deskElem.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("desk's elem expect struct, got %s", deskElem.Elem().Kind())
		}

		rv := reflect.ValueOf(input)
		rt := reflect.TypeOf(input)
		elemType := rt.Elem()

		var indexVal, columnVal reflect.Value
		direct := reflect.Indirect(deskValue)
		mapReflect := reflect.MakeMap(deskElem)
		deskKey := deskValue.Type().Elem().Key()
		for i := 0; i < rv.Len(); i++ {
			curVal := rv.Index(i)
			indexVal, columnVal, err = findStructValByIndexKey(curVal, elemType, indexKey, columnKey)

			if err != nil {
				return
			}
			if deskKey.Kind() != indexVal.Kind() {
				return fmt.Errorf("cant't convert %s to %s, your map'key must be %s", indexVal.Kind(), deskKey.Kind(), indexVal.Kind())
			}
			if len(columnKey) == 0 {
				mapReflect.SetMapIndex(indexVal, curVal)
				direct.Set(mapReflect)
			} else {
				if deskElem.Elem().Kind() != columnVal.Kind() {
					return fmt.Errorf("your map must be map[%s]%s", indexVal.Kind(), columnVal.Kind())
				}
				mapReflect.SetMapIndex(indexVal, columnVal)
				direct.Set(mapReflect)
			}
		}
		return
	}

	structColumn := func(desk, input interface{}, columnKey string) (err error) {
		findStructValByColumnKey := func(curVal reflect.Value, elemType reflect.Type, columnKey string) (columnVal reflect.Value, err error) {
			columnExist := false
			for i := 0; i < elemType.NumField(); i++ {
				curField := curVal.Field(i)
				if elemType.Field(i).Name == columnKey {
					columnExist = true
					columnVal = curField
					continue
				}
			}
			if !columnExist {
				return columnVal, fmt.Errorf("columnKey %s not found in %s's field", columnKey, elemType)
			}
			return
		}

		if len(columnKey) == 0 {
			return fmt.Errorf("columnKey cannot not be empty")
		}

		deskElemType := reflect.TypeOf(desk).Elem()
		if deskElemType.Kind() != reflect.Slice {
			return fmt.Errorf("desk must be slice")
		}

		rv := reflect.ValueOf(input)
		rt := reflect.TypeOf(input)

		var columnVal reflect.Value
		deskValue := reflect.ValueOf(desk)
		direct := reflect.Indirect(deskValue)

		for i := 0; i < rv.Len(); i++ {
			columnVal, err = findStructValByColumnKey(rv.Index(i), rt.Elem(), columnKey)
			if err != nil {
				return
			}
			if deskElemType.Elem().Kind() != columnVal.Kind() {
				return fmt.Errorf("your slice must be []%s", columnVal.Kind())
			}

			direct.Set(reflect.Append(direct, columnVal))
		}
		return
	}

	deskValue := reflect.ValueOf(desk)
	if deskValue.Kind() != reflect.Ptr {
		return fmt.Errorf("desk must be ptr")
	}

	rv := reflect.ValueOf(input)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return fmt.Errorf("input must be map slice or array")
	}

	rt := reflect.TypeOf(input)
	if rt.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("input's elem must be struct")
	}

	if len(indexKey) > 0 {
		return structIndexColumn(desk, input, columnKey, indexKey)
	}
	return structColumn(desk, input, columnKey)
}

/**
 * @Author: mali
 * @Func:
 * @Description: 结构体转map
 * @Param:
 * @Return:
 * @Example:
 * @param {interface{}} in
 * @param {string} tagName 使用结构体里面指定的tag标签当map对应的key
 * @param {string} index tag标签,分割的第几个值为key
 */
func StructToMap(in interface{}, tagName string, index ...int) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("非结构体或结构体指针不适配; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			tag_calue_list := strings.Split(tagValue, ",")
			if len(index) > 0 {
				out[tag_calue_list[index[0]]] = v.Field(i).Interface()
			} else {
				out[tag_calue_list[0]] = v.Field(i).Interface()
			}

		}
	}
	return out, nil
}

/**
 * @Author: mali
 * @Func:
 * @Description: 结构体转map 适配无限级匿名结构体字段嵌套的转化
 * @Param:
 * @Return:
 * @Example:
 * @param {interface{}} in
 * @param {string} tagName 使用结构体里面指定的tag标签当map对应的key
 * @param {string} index tag标签,分割的第几个值为key
 */
func NewStructToMap(in interface{}, tagName string, index ...int) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("非结构体或结构体指针不适配; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		//匿名结构体字段
		if fi.Anonymous && fi.Type.Kind() == reflect.Struct {
			res, _ := NewStructToMap(v.Field(i).Interface(), tagName, index...)
			for k, v := range res {
				if _, ok := out[k]; !ok {
					out[k] = v
				}
			}
		} else {
			if tagValue := fi.Tag.Get(tagName); tagValue != "" {
				tag_calue_list := strings.Split(tagValue, ",")
				if len(index) > 0 {
					out[tag_calue_list[index[0]]] = v.Field(i).Interface()
				} else {
					out[tag_calue_list[0]] = v.Field(i).Interface()
				}

			}
		}

	}
	return out, nil
}
