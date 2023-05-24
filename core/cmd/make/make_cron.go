/*
 * @Description:
 * @Author: mali
 * @Date: 2023-04-14 14:20:37
 * @LastEditTime: 2023-05-24 16:40:09
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

var CmdMakeCron = &cobra.Command{
	Use:   "cron",
	Short: "生成自定义定时任务的指令，示例: make cron test",
	Run:   runMakeCron,
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func runMakeCron(cmd *cobra.Command, args []string) {
	structName := strcase.ToCamel(args[0])
	model := Model{
		StructName:  structName,
		PackageName: str.LowerCamel(args[0]),
	}
	// 拼接目标文件路径
	filePath := fmt.Sprintf("app/cron/cron_job/%s.go", model.PackageName)

	// 从模板中创建文件（做好变量替换）
	createFileFromStub(filePath, "cron", model)

	// 友好提示
	console.Success("指令文件已经生成在app/cron/cron_job 下，请自行修改 指令名为:" + model.PackageName)
}
