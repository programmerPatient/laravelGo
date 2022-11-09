/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-20 16:41:02
 * @LastEditTime: 2022-09-20 16:50:21
 * @LastEditors: VSCode
 * @Reference:
 */
package job

import (
	"context"

	"github.com/hibiken/asynq"
)

var JobNameToHandle = map[string]func(ctx context.Context, t *asynq.Task) error{}

func AddJobNameToHandle(name string, handler func(ctx context.Context, t *asynq.Task) error) {
	JobNameToHandle[name] = handler
}
