/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-20 15:30:57
 * @LastEditTime: 2022-09-20 15:35:19
 * @LastEditors: VSCode
 * @Reference:
 */
package queue

import (
	"github.com/laravelGo/core/asynq"
	"github.com/spf13/cobra"
)

var CmdQueueServer = &cobra.Command{
	Use:   "server",
	Short: "执行队列任务",
	Run:   RunQueueServer,
}

func RunQueueServer(cmd *cobra.Command, args []string) {
	asynq.Exec()
}
