/*
 * @Description:随机
 * @Author: mali
 * @Date: 2023-07-31 15:13:20
 * @LastEditTime: 2023-07-31 15:13:24
 * @LastEditors: VSCode
 * @Reference:
 */
package helper

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cast"
)

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
 * @Description: 生成区间随机数
 * @Param:
 * @Return:
 * @Example:
 * @param {*} min
 * @param {int64} max
 */
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

/**
 * @Author: mali
 * @Func:
 * @Description: 随机切片元素
 * @Param:
 * @Return:
 * @Example:
 * @param {[]interface{}} elements
 */
func ArrayRand(input []interface{}, length int) []interface{} {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	input_length := len(input)
	if length > input_length {
		length = input_length
	}
	n := make([]interface{}, input_length)
	for i, v := range r.Perm(input_length) {
		n[i] = input[v]
	}
	return_data := make([]interface{}, length)
	for i := 0; i < length; i++ {
		return_data[i] = n[i]
	}
	return return_data
}
