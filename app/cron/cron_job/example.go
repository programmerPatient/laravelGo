/*
 * @Description:
 * @Author: mali
 * @Date: 2023-10-09 13:34:45
 * @LastEditTime: 2023-10-09 13:35:19
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

type Example struct {
}

var CronExample = &Example{}

//！！！！!!记得生成的指令要加入到app/cron/kernel.go 下的InintCron里面才能生效!!!!!!!!!

// 定时任务名
func (c *Example) GetCronName() string {
	return "example"
}

// 默认每次启动命令是立即开始执行的函数
func (c *Example) GetStartDefaultRunFunc() {

}

// 初始化配置信息
func (c *Example) InitServer() {
}

// 错误捕获防止当该定时任务报错时，影响其他定时任务的执行
func (c *Example) PanicRecover(err interface{}) {
	//添加自己的错误处理逻辑
	s := string(debug.Stack())
	fmt.Printf("err=%v, stack=%s\n", err, s)
}

// 定时任务执行规则，自行修改
func (c *Example) GetSpec() string {
	return "*/2 * * * * ?" //每隔 2 秒钟 执行一次任务
}

// 定时任务执行回调函数，替换成自己的逻辑
func (c *Example) Run() {
	console.Success("定时任务执行开始！！！")
	console.Success(fmt.Sprintf("%v任务执行时间为：%v", c.GetCronName(), time.Now().Format("2006-01-02 15:04:05")))
	//pass
	console.Success("定时任务执行结束！！！")
}
