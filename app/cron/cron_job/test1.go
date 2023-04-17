package cron_job

import "github.com/laravelGo/core/console"
type Test1 struct {
}
var CronTest1 = &Test1{}

//！！！！!!记得生成的指令要加入到app/cron/kernel.go 下的InintCron里面才能生效!!!!!!!!!

//定时任务名
func (c *Test1) GetCronName() string {
	return "test1"
}
//默认每次启动命令是立即开始执行的函数
func (c *Test1) GetStartDefaultRunFunc() func() {
	return func() {
		
	}
}

//定时任务执行规则，自行修改
func (c *Test1) GetSpec() string {
	return "*/2 * * * * ?" //每隔 2 秒钟 执行一次任务
}

//定时任务执行回调函数，替换成自己的逻辑
func (c *Test1) Run() func()  {
	return func(){
		console.Success("定时任务执行开始！！！")
		console.Success("这是一个定时任务")
		//pass 
		console.Success("定时任务执行结束！！！")
	}
}
