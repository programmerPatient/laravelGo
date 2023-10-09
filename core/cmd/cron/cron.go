/*
 * @Description:定时任务
 * @Author: mali
 * @Date: 2023-04-14 14:19:04
 * @LastEditTime: 2023-10-09 13:32:04
 * @LastEditors: VSCode
 * @Reference:
 */
package cron

import (
	"fmt"
	"sync"

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
	var id_map_name sync.Map          //缓存每个定时任务对应的id,防止并发读写问题
	for _, v := range workcron {
		value := v.(BaseCron)
		if len(args) == 0 || helper.InArray(value.GetCronName(), args) {
			console.Success(fmt.Sprintf("定时任务：【%v】开始运行", value.GetCronName()))
			value.InitServer()
			go value.GetStartDefaultRunFunc()
			cron_id, err := c.AddFunc(value.GetSpec(), func() {
				defer func() {
					if err := recover(); err != nil {
						console.Error(fmt.Sprintf("定时任务【%v】运行错误:%v", value.GetCronName(), err))
						value.PanicRecover(err)
						//移除定时任务
						id, ok := id_map_name.Load(value.GetCronName())
						if ok {
							c.Remove(id.(cron.EntryID))
						}
					}
				}()
				value.Run()
			})
			if err != nil {
				console.Error(fmt.Sprintf("定时任务【%v】添加失败:%v", value.GetCronName(), err))
				continue
			}
			id_map_name.Store(value.GetCronName(), cron_id)
		}
	}
	c.Start()
	select {}
}
