/*
 * @Description:生成job
 * @Author: mali
 * @Date: 2022-09-21 09:39:38
 * @LastEditTime: 2022-09-21 09:50:25
 * @LastEditors: VSCode
 * @Reference:
 */
package make

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/laravelGo/core/console"
	"github.com/laravelGo/core/str"
	"github.com/spf13/cobra"
)

var CmdMakeJob = &cobra.Command{
	Use:   "job",
	Short: "创建job",
	Run:   runMakeJob,
	Args:  cobra.ExactArgs(1), //只有一个参数
}

func runMakeJob(cmd *cobra.Command, args []string) {
	console.Warning("开始生成任务文件")
	//文件名称
	fileName := str.Snake(args[0])
	structName := strcase.ToCamel(args[0])
	filePath := fmt.Sprintf("app/job/%s.go", fileName)
	model := &Model{
		StructName: structName,
		TableName:  str.LowerCamel(args[0]),
	}
	createFileFromStub(filePath, "job", *model)
	console.Success("任务文件已经生成")
}
