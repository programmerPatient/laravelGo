/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-20 16:14:03
 * @LastEditTime: 2022-09-20 16:56:13
 * @LastEditors: VSCode
 * @Reference:
 */
package asynq

import (
	"context"

	"github.com/hibiken/asynq"
)

type AsynqJob interface {
	Handle() func(context.Context, *asynq.Task) error
	GetName() string
	GetPayload() interface{}
	GetOpt() []asynq.Option
}
