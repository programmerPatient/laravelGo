/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-08 17:02:05
 * @LastEditTime: 2022-11-08 11:30:43
 * @LastEditors: VSCode
 * @Reference:
 */
package helper

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// MicrosecondsStr 将 time.Duration 类型（nano seconds 为单位）
// 输出为小数点后 3 位的 ms （microsecond 毫秒，千分之一秒）
func MicrosecondsStr(elapsed time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
}

func ArrayColumn(input map[string]map[string]interface{}, columnKey string) []interface{} {
	columns := make([]interface{}, 0, len(input))
	for _, val := range input {
		if v, ok := val[columnKey]; ok {
			columns = append(columns, v)
		}
	}
	return columns
}

func InArray(needle interface{}, hystack interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range hystack.([]string) {
			if key == item {
				return true
			}
		}
	case uint:
		for _, item := range hystack.([]uint) {
			if key == item {
				return true
			}
		}
	case uint8:
		for _, item := range hystack.([]uint8) {
			if key == item {
				return true
			}
		}
	case uint16:
		for _, item := range hystack.([]uint16) {
			if key == item {
				return true
			}
		}
	case uint32:
		for _, item := range hystack.([]uint32) {
			if key == item {
				return true
			}
		}
	case uint64:
		for _, item := range hystack.([]uint64) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range hystack.([]int) {
			if key == item {
				return true
			}
		}
	case int8:
		for _, item := range hystack.([]int8) {
			if key == item {
				return true
			}
		}
	case int16:
		for _, item := range hystack.([]int16) {
			if key == item {
				return true
			}
		}
	case int32:
		for _, item := range hystack.([]int32) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range hystack.([]int64) {
			if key == item {
				return true
			}
		}
	case float32:
		for _, item := range hystack.([]float32) {
			if key == item {
				return true
			}
		}
	case float64:
		for _, item := range hystack.([]float64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
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
 * @Description: Empty 类似于 PHP 的 empty() 函数 判断是否为空元素
 * @Param:
 * @Return:
 * @Example:
 * @param {interface{}} val
 */
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

func Http(method string, url string, data interface{}, headers ...map[string]string) (string, error) {
	client := &http.Client{}
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Length", cast.ToString(req.ContentLength))
	if len(headers) > 0 {
		for key, header := range headers[0] {
			req.Header.Set(key, header)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//关键的一步，去除最后的空格防止读取报错invalid header field name
	body_string := strings.TrimSpace(string(body))
	return body_string, nil
}

func HttpsPostForm(url string, data url.Values) (string, error) {
	resp, err := http.PostForm(url, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

/**
 * @Author: mali
 * @Func:
 * @Description: 秒数转化为分秒字符串
 * @Param:
 * @Return:
 * @Example:
 * @param {int64} se
 */
func SecondToMsString(se int64) string {
	if se > 60 {
		min := se / 60
		second := se % 60
		return fmt.Sprintf("%v分%v秒", min, second)
	} else {
		return fmt.Sprintf("%v秒", se)
	}
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

func HttpBuildQuery(params map[string]string) (param_str string) {
	params_arr := make([]string, 0, len(params))
	for k, v := range params {
		params_arr = append(params_arr, fmt.Sprintf("%s=%s", k, v))
	}
	param_str = strings.Join(params_arr, "&")
	return param_str
}

func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has) //将[]byte转成16进制
}

/**
从url中获取图片资源
*/
func ReadImgData(url string) image.Image {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil
	}
	return img
}

//保存图片
func SaveImage(targetPath string, m image.Image) error {
	fSave, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer fSave.Close()

	err = jpeg.Encode(fSave, m, nil)

	if err != nil {
		return err
	}

	return nil
}

//PostJson 获取post json参数
func PostJson(ctx *gin.Context, obj interface{}) error {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}
	return nil
}

/**
 * @Author: mali
 * @Func:
 * @Description: 随机长度的数字字符串
 * @Param:
 * @Return:
 * @Example:
 * @param {string} length
 */
func RandString(length int) string {
	var n int32 = 1
	for i := 0; i < length; i++ {
		n *= 10
	}
	return fmt.Sprintf("%0"+cast.ToString(length)+"v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(n))
}

/**
  * @Author: mali
  * @Func:
  * @Description: stirng map合并
  * @Param:
  * @Return:
  * @Example:
  *      {
			 a:map[string]string{"1": "110", "2":"120", "3":"119"}
			 b:map[string]string{"1": "111", "2":"122", "4":"129"}
			 return:map["1":111 "2":122 "3":119 "4":129]
		  }
  * @param {map[string]interface{}} a
  * @param {map[string]interface{}} b
*/
func StringMapMerge(a map[string]interface{}, b map[string]interface{}) map[string]interface{} {
	n := make(map[string]interface{})
	for i, v := range a {
		for j, w := range b {
			if i == j {
				n[i] = w
			} else {
				if _, ok := n[i]; !ok {
					n[i] = v
				}
				if _, ok := n[j]; !ok {
					n[j] = w
				}
			}
		}
	}
	return n
}

/**
 * @Author: mali
 * @Func:
 * @Description: 判断文件或文件夹是否存在
 * @Param:
 * @Return:
 * @Example:
 * @param {string} filename
 */
func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func Sha1(s string) string {
	//产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，然后 sha1.Sum([]byte{})。这里我们从一个新的散列开始。
	h := sha1.New() // md5加密类似md5.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write([]byte(s))
	//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来对现有的字符切片追加额外的字节切片：一般不需要要。
	bs := h.Sum(nil)
	//SHA1 值经常以 16 进制输出，使用%x 来将散列结果格式化为 16 进制字符串。
	return fmt.Sprintf("%x", bs)
}
