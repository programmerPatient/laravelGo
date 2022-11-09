/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-20 15:31:01
 * @LastEditTime: 2022-09-20 15:35:10
 * @LastEditors: VSCode
 * @Reference:
 */
package queue

import "github.com/spf13/cobra"

var CmdQueue = &cobra.Command{
	Use:   "queue",
	Short: "队列相关指令",
}

func init() {
	// 注册 queue 的子命令
	CmdQueue.AddCommand(
		CmdQueueServer,
	)
}
