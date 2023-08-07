/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-08 17:02:05
 * @LastEditTime: 2023-07-31 15:16:55
 * @LastEditors: VSCode
 * @Reference:
 */
package helper

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// MicrosecondsStr 将 time.Duration 类型（nano seconds 为单位）
// 输出为小数点后 3 位的 ms （microsecond 毫秒，千分之一秒）
func MicrosecondsStr(elapsed time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
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

// 下载图片信息
func DownLoadImage(url string, base string) (string, error) {

	pic := base
	idx := strings.LastIndex(url, "/")
	if idx < 0 {
		pic += "/" + url
	} else {
		pic += url[idx+1:]
	}
	v, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer v.Body.Close()
	content, err := ioutil.ReadAll(v.Body)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(pic, content, 0666)
	if err != nil {
		return "", err
	}
	return pic, nil
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

//PostJson 获取post 固定key的参数json参数
func PostJsonOnly(ctx *gin.Context, obj *map[string]interface{}, field []string) error {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	res := map[string]interface{}{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}
	for k, v := range res {
		if InArray(k, field) {
			(*obj)[k] = v
		}
	}
	return nil
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

/**
 * @Author: mali
 * @Func:
 * @Description: 十进制转化为2、8、16进制
 * @Param:
 * @Return:
 * @Example:
 * @param {*} n
 * @param {int} num
 */
func DecConvertToX(n int, num int) (string, error) {
	if n < 0 {
		return strconv.Itoa(n), errors.New("只支持正整数")
	}
	if num != 2 && num != 8 && num != 16 {
		return strconv.Itoa(n), errors.New("只支持二、八、十六进制的转换")
	}
	result := ""
	h := map[int]string{
		0:  "0",
		1:  "1",
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "A",
		11: "B",
		12: "C",
		13: "D",
		14: "E",
		15: "F",
	}
	for ; n > 0; n /= num {
		lsb := h[n%num]
		result = lsb + result
	}
	return result, nil
}

/**
 * @Author: mali
 * @Func:
 * @Description: 获取唯一的id
 * @Param:
 * @Return:
 * @Example:
 */
func GetUuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

/**
 * @Author: mali
 * @Func:
 * @Description: 四舍五入保留固定位子小数
 * @Param:
 * @Return:
 * @Example:
 * @param {float64} val
 * @param {int} precision
 */
func Round(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Floor(val*p+0.5) / p
}

// CheckSignature 微信公众号签名检查
func CheckSignature(signature, timestamp, nonce, token string) bool {
	arr := []string{timestamp, nonce, token}
	// 字典序排序
	sort.Strings(arr)

	n := len(timestamp) + len(nonce) + len(token)
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < len(arr); i++ {
		b.WriteString(arr[i])
	}

	return Sha1(b.String()) == signature
}

//进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

/**
 * @Author: mali
 * @Func:
 * @Description: 判断是否为数字类型
 * @Param:
 * @Return:
 * @Example:
 * @param {interface{}} val
 */
func IsNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if str == "" {
			return false
		}
		// Trim any whitespace
		str = strings.TrimSpace(str)
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		// hex
		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9,Point,Scientific
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' { // Point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}

	return false
}
