/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-20 15:15:01
 * @LastEditTime: 2022-09-21 15:18:17
 * @LastEditors: VSCode
 * @Reference:
 */
package job

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/laravelGo/core/logger"
)

//任务名
const ExampleJobName = "test"

type ExampleJob struct {
	Name     string
	Opts     []asynq.Option
	Playload interface{}
}

//参数结构体
type ExamplePayload struct {
	UserId int
}

func init() {
	AddJobNameToHandle(ExampleJobName, (&ExampleJob{}).Handle())
}

//新建任务
func NewExampletJob(playload interface{}) *ExampleJob {
	return &ExampleJob{
		Name:     ExampleJobName,
		Playload: playload,
		Opts: []asynq.Option{
			asynq.Queue("critical"),
			asynq.ProcessIn(10 * time.Second),
		},
	}
}

//任务名
func (job *ExampleJob) GetName() string {
	return ExampleJobName
}

//参数
func (job *ExampleJob) GetPayload() interface{} {
	return job.Playload
}

//任务处理行为
func (job *ExampleJob) GetOpt() []asynq.Option {
	return job.Opts
}

//处理函数
func (job *ExampleJob) Handle() func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var p ExamplePayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			logger.Dump(err.Error())
		}
		logger.Dump(p)
		return nil
	}
}
