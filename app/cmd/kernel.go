/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-05 15:06:56
 * @LastEditTime: 2022-09-20 13:56:12
 * @LastEditors: VSCode
 * @Reference:
 */
package cmd

import (
	"github.com/laravelGo/app/cmd/commands"
	"github.com/spf13/cobra"
)

//Cmds 放置自定义命令行指令
var Cmds []*cobra.Command

/**
 * @Author: mali
 * @Func: InintCmd
 * @Description: 加载用户自定义cmd指令
 * @Param:
 * @Return:
 * @Example:
 */
//
func InintCmd() []*cobra.Command {
	Cmds = append(Cmds,
		commands.CmdExample,
	)
	return Cmds
}
