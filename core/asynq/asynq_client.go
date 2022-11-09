/*
 * @Description:异步任务客户端
 * @Author: mali
 * @Date: 2022-09-20 14:27:03
 * @LastEditTime: 2022-09-21 14:20:38
 * @LastEditors: VSCode
 * @Reference:
 */
package asynq

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hibiken/asynq"
	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/logger"
	"go.uber.org/zap"
)

var AsynqClient *asynq.Client

var once sync.Once

func GetAsynqClient() *asynq.Client {
	once.Do(func() {
		AsynqClient = asynq.NewClient(
			asynq.RedisClientOpt{
				Addr:     config.GetString("asynq.redis.host") + ":" + config.GetString("asynq.redis.port"),
				Username: config.GetString("asynq.redis.username"),
				Password: config.GetString("asynq.redis.password"),
				DB:       config.GetInt("asynq.redis.database"),
			},
		)
	})
	return AsynqClient
}

/**
 * @Author: mali
 * @Func:
 * @Description: 任务投递
 * @Param:
 * @Return:
 * @Example:
 * @param {string} typename
 * @param {interface{}} payload
 * @param {...asynq.Option} opts
 */
func NewTask(typename string, payload interface{}, opts ...asynq.Option) (*asynq.Task, error) {
	payload_byte, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(typename, payload_byte, opts...), nil
}

/**
 * @Author: mali
 * @Func:
 * @Description: 投递任务
 * @Param:
 * @Return:
 * @Example:
 * @param {string} typename
 * @param {interface{}} payload
 * @param {...asynq.Option} opts
 */
func Delivery(job AsynqJob) {
	task, err := NewTask(job.GetName(), job.GetPayload(), job.GetOpt()...)
	if err != nil {
		logger.Error("asynq-delivery", zap.String(job.GetName(), err.Error()))
	}
	// Process the task immediately.
	info, err := GetAsynqClient().Enqueue(task, job.GetOpt()...)
	if err != nil {
		logger.Error("asynq-delivery", zap.String(job.GetName(), err.Error()))
	} else {
		logger.Debug(fmt.Sprintf(" [*] Successfully enqueued task: %+v", info))
	}

}
