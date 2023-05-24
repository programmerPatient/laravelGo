/*
 * @Description:
 * @Author: mali
 * @Date: 2023-04-14 14:32:38
 * @LastEditTime: 2023-05-24 16:32:09
 * @LastEditors: VSCode
 * @Reference:
 */
package cron

type BaseCron interface {
	GetCronName() string            //定时任务名
	GetSpec() string                //定时任务执行规则
	GetStartDefaultRunFunc() func() //默认每次启动命令是立即开始执行的函数
	PanicRecover()                  //如果程序发现问题，可以使用此函数来捕捉并根据实际情况重
	Run() func()                    //执行任务函数
}
