/*
 * @Description:
 * @Author: mali打印数据
 * @Date: 2023-08-04 14:24:28
 * @LastEditTime: 2023-08-04 15:02:30
 * @LastEditors: VSCode
 * @Reference:
 */
package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
)

/**
 * @Author: mali
 * @Func:
 * @Description: 格式化打印数据，提升数据打印可读性
 * @Param:
 * @Return:
 * @Example:
 * @param {interface{}} data
 */
func FormatPrint(data interface{}) {
	bs, _ := json.Marshal(data)
	var out bytes.Buffer
	json.Indent(&out, bs, "", "  ")
	fmt.Println("\n", out.String())
}
