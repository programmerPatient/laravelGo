/*
 * @Description:异步队列配置
 * @Author: mali
 * @Date: 2022-09-21 14:06:06
 * @LastEditTime: 2022-09-21 15:21:02
 * @LastEditors: VSCode
 * @Reference:
 */
package config

import "github.com/laravelGo/core/config"

func init() {
	config.AddConfig("asynq", func() map[string]interface{} {
		return map[string]interface{}{
			// 最大并发处理任务数
			// 如果设置为零或负值，则为 CPU 数量
			"concurrency": 10,
			// redis配置 默认使用env里面的配置，也可单独设置
			"redis": map[string]interface{}{
				"host":     config.Env("REDIS_HOST", "127.0.0.1"),
				"port":     config.Env("REDIS_PORT", "6379"),
				"password": config.Env("REDIS_PASSWORD", ""),
				"database": config.Env("REDIS_MAIN_DB", 1),
			},
			//StrictPriority 表示是否应该严格对待队列优先级
			//如果设置为true，则优先处理队列中优先级最高的任务
			//只有当优先级高的队列为空时，才会处理低优先级队列中的任务
			"strict_priority": true,
			//不开启strict_priority时
			//队列的优先级 此处代表三个队列 与队列名称关联的数字是队列的优先级
			//critical任务将被处理 3/6 的时间
			//default任务将被处理 2/6 的时间
			//low任务将被处理 1/6 的时间
			"queue": map[string]interface{}{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		}
	})
}
