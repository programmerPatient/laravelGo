/*
 * @Description:异步任务服务端
 * @Author: mali
 * @Date: 2022-09-20 14:53:21
 * @LastEditTime: 2022-09-21 15:00:38
 * @LastEditors: VSCode
 * @Reference:
 */
package asynq

import (
	"github.com/hibiken/asynq"
	"github.com/laravelGo/app/job"
	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/logger"
	"go.uber.org/zap"
)

func Exec() {
	AsynqServer := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     config.GetString("asynq.redis.host") + ":" + config.GetString("asynq.redis.port"),
			Username: config.GetString("asynq.redis.username"),
			Password: config.GetString("asynq.redis.password"),
			DB:       config.GetInt("asynq.redis.database"),
		},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: config.GetInt("asynq.concurrency"),
			// Optionally specify multiple queues with different priority.
			Queues:         config.GetStringMapInt("asynq.queue"),
			StrictPriority: config.GetBool("asynq.strict_priority"),
			// See the godoc for other configuration options
		},
	)
	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	for type_name, handler := range job.JobNameToHandle {
		mux.HandleFunc(type_name, handler)
	}
	if err := AsynqServer.Run(mux); err != nil {
		logger.Error("asynq", zap.String("error", err.Error()))
	}
}
