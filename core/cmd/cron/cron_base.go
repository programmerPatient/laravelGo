/*
 * @Description:
 * @Author: mali
 * @Date: 2023-04-14 14:32:38
 * @LastEditTime: 2023-10-09 10:53:33
 * @LastEditors: VSCode
 * @Reference:
 */
package cron

type BaseCron interface {
	GetCronName() string            //定时任务名
	GetSpec() string                //定时任务执行规则
	GetStartDefaultRunFunc() func() //默认每次启动命令是立即开始执行的函数
	InitServer()                    //初始化一些额外需要的额服务
	Run() func()                    //执行任务函数
	PanicRecover(err interface{})   //异常处理函数
}
