/*
 * @Description:生成cmd指令模版
 * @Author: mali
 * @Date: 2022-09-14 10:43:26
 * @LastEditTime: 2022-09-21 09:52:15
 * @LastEditors: VSCode
 * @Reference:
 */
package make

import (
	"fmt"

	"github.com/laravelGo/core/console"

	"github.com/spf13/cobra"
)

var CmdMakeCMD = &cobra.Command{
	Use:   "cmd",
	Short: "生成自定义指令的指令，示例: make cmd test",
	Run:   runMakeCMD,
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func runMakeCMD(cmd *cobra.Command, args []string) {

	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(args[0])

	// 拼接目标文件路径
	filePath := fmt.Sprintf("app/cmd/commands/%s.go", model.PackageName)

	// 从模板中创建文件（做好变量替换）
	createFileFromStub(filePath, "cmd", model)

	// 友好提示
	console.Success("指令文件已经生成在app/cmd/commands/ 下，请自行修改 指令名为:" + model.PackageName)
}
