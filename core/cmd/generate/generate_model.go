/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-15 17:50:28
 * @LastEditTime: 2022-09-15 18:00:02
 * @LastEditors: VSCode
 * @Reference:
 */
package generate

import (
	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/generate"
	"github.com/spf13/cobra"
)

var CmdGenerateModel = &cobra.Command{
	Use:   "model",
	Short: "根据数据库表生成模型结构体 不传参数则库总所有的表都会生成 传参数代表生成指定的表",
	Args:  cobra.MinimumNArgs(0),
	Run:   runGengrate,
}

func runGengrate(cmd *cobra.Command, args []string) {
	generate.Generate(config.GetString("database.mysql.database"), args...)
}
