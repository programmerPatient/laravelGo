/*
 * @Description:
 * @Author: mali
 * @Date: 2023-04-17 09:20:07
 * @LastEditTime: 2023-10-09 11:25:37
 * @LastEditors: VSCode
 * @Reference:
 */
package cron_job

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/laravelGo/core/console"
)

type Test struct {
}

var CronTest = &Test{}

//！！！！!!记得生成的指令要加入到app/cron/kernel.go 下的InintCron里面才能生效!!!!!!!!!

// 定时任务名
func (c *Test) GetCronName() string {
	return "test"
}

// 初始化配置信息
func (c *Test) InitServer() {
}

// 定时任务执行描述
func (c *Test) GetCronDescription() string {
	return "添加任务描述"
}

// 默认每次启动命令是立即开始执行的函数
func (c *Test) GetStartDefaultRunFunc() func() {
	return func() {}
}

// 定时任务执行规则，自行修改
func (c *Test) GetSpec() string {
	return "*/2 * * * * ?" //每隔 2 秒钟 执行一次任务
}

// 错误捕获防止当该定时任务报错时，影响其他定时任务的执行
func (c *Test) PanicRecover(err interface{}) {
	//添加自己的错误处理逻辑
	s := string(debug.Stack())
	fmt.Printf("err=%v, stack=%s\n", err, s)
}

// 定时任务执行回调函数，替换成自己的逻辑
func (c *Test) Run() func() {
	return func() {
		panic("dsad")
		//错误捕获防止当该定时任务报错时，中断所有定时任务
		console.Success(fmt.Sprintf("%v任务执行时间为：%v", c.GetCronName(), time.Now().Format("2006-01-02 15:04:05")))
	}
}
