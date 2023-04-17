/*
 * @Description:
 * @Author: mali
 * @Date: 2023-04-17 09:20:07
 * @LastEditTime: 2023-04-17 09:20:14
 * @LastEditors: VSCode
 * @Reference:
 */
package cron_job

import (
	"fmt"
	"time"

	"github.com/laravelGo/core/console"
)

type Test struct {
}

var CronTest = &Test{}

//！！！！!!记得生成的指令要加入到app/cron/kernel.go 下的InintCron里面才能生效!!!!!!!!!

//定时任务名
func (c *Test) GetCronName() string {
	return "test"
}

//定时任务执行描述
func (c *Test) GetCronDescription() string {
	return "添加任务描述"
}

//默认每次启动命令是立即开始执行的函数
func (c *Test) GetStartDefaultRunFunc() func() {
	return func() {}
}

//定时任务执行规则，自行修改
func (c *Test) GetSpec() string {
	return "*/2 * * * * ?" //每隔 2 秒钟 执行一次任务
}

//定时任务执行回调函数，替换成自己的逻辑
func (c *Test) Run() func() {
	return func() {
		console.Success(fmt.Sprintf("%v任务执行时间为：%v", c.GetCronName(), time.Now().Format("2006-01-02 15:04:05")))
	}
}
