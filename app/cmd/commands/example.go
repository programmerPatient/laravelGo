/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-05 15:15:00
 * @LastEditTime: 2022-09-20 13:56:07
 * @LastEditors: VSCode
 * @Reference:
 */
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var CmdExample = &cobra.Command{
	Use:   "example",
	Short: "example 示例指令",
	Args:  cobra.NoArgs, //不需要参数
	Run:   runHello,
}

func runHello(cmd *cobra.Command, args []string) {
	fmt.Println("hello word")
}
