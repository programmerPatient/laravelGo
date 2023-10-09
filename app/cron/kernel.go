/*
 * @Description:
 * @Author: mali
 * @Date: 2023-04-17 09:18:48
 * @LastEditTime: 2023-10-09 13:36:21
 * @LastEditors: VSCode
 * @Reference:
 */
package cron

import "github.com/laravelGo/app/cron/cron_job"

// 加载用户自定义定时任务，不在这里面的定时任务不会被执行
func InintCron() []interface{} {
	return []interface{}{
		cron_job.CronExample,
	}
}
