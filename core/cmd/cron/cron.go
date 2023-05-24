/*
 * @Description:定时任务
 * @Author: mali
 * @Date: 2023-04-14 14:19:04
 * @LastEditTime: 2023-05-24 16:38:15
 * @LastEditors: VSCode
 * @Reference:
 */
package cron

import (
	"fmt"

	work_cron "github.com/laravelGo/app/cron"
	"github.com/laravelGo/app/helper"
	"github.com/laravelGo/core/console"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var CronCmd = &cobra.Command{
	Use:   "cron",
	Short: "开启定时任务,可以指定只执行固定的定时任务,不传递参数为执行全部的定时任务 example:go run main.go cron cron1 cron2 cron3",
	Args:  cobra.MinimumNArgs(0),
	Run:   runCron,
}

func runCron(cmd *cobra.Command, args []string) {
	workcron := work_cron.InintCron()
	c := cron.New(cron.WithSeconds()) //精确到秒
	for _, v := range workcron {
		value := v.(BaseCron)
		if len(args) == 0 || helper.InArray(value.GetCronName(), args) {
			console.Success(fmt.Sprintf("定时任务：【%v】开始运行", value.GetCronName()))
			//异步执行任务
			go value.GetStartDefaultRunFunc()()

			c.AddFunc(value.GetSpec(), value.Run())
		}
	}
	c.Start()
	select {}
}
