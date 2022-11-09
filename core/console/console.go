/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-05 14:24:55
 * @LastEditTime: 2022-09-05 14:26:06
 * @LastEditors: VSCode
 * @Reference:
 */
// Package console 命令行辅助方法
package console

import (
	"fmt"
	"os"

	"github.com/mgutz/ansi"
)

/**
 * @Author: mali
 * @Func: Success
 * @Description: 打印一条成功消息，绿色输出
 * @Param:
 * @Return:
 * @Example:
 * @param {string} msg
 */
func Success(msg string) {
	colorOut(msg, "green")
}

/**
 * @Author: mali
 * @Func: Error
 * @Description: 打印一条报错消息，红色输出
 * @Param:
 * @Return:
 * @Example:
 * @param {string} msg
 */
func Error(msg string) {
	colorOut(msg, "red")
}

/**
 * @Author: mali
 * @Func: Warning
 * @Description: 打印一条提示消息，黄色输出
 * @Param:
 * @Return:
 * @Example:
 * @param {string} msg
 */
func Warning(msg string) {
	colorOut(msg, "yellow")
}

/**
 * @Author: mali
 * @Func: Exit
 * @Description: 打印一条报错消息，并退出 os.Exit(1)
 * @Param:
 * @Return:
 * @Example:
 * @param {string} msg
 */
func Exit(msg string) {
	Error(msg)
	os.Exit(1)
}

/**
 * @Author: mali
 * @Func: ExitIf
 * @Description: 语法糖，自带 err != nil 判断
 * @Param:
 * @Return:
 * @Example:
 * @param {error} err
 */
func ExitIf(err error) {
	if err != nil {
		Exit(err.Error())
	}
}

/**
 * @Author: mali
 * @Func: colorOut
 * @Description: 内部使用，设置高亮颜色
 * @Param:
 * @Return:
 * @Example:
 * @param {*} message
 * @param {string} color
 */
func colorOut(message, color string) {
	fmt.Fprintln(os.Stdout, ansi.Color(message, color))
}
