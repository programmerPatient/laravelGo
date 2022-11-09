/*
 * @Author: error: git config user.name && git config user.email & please set dead value or install git
 * @Date: 2022-08-25 13:56:36
 * @LastEditors: VSCode
 * @LastEditTime: 2022-09-15 18:00:32
 * @FilePath: /laravelGo/laravego/cmd/generate/generate.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package generate

import (
	"github.com/spf13/cobra"
)

var CmdGenerate = &cobra.Command{
	Use:   "generate",
	Short: "自动生成命令",
}

func init() {
	CmdGenerate.AddCommand(
		CmdGenerateKey,
		CmdGenerateModel,
	)
}
