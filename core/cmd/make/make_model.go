/*
 * @Description:生成模型文件指令
 * @Author: mali
 * @Date: 2022-09-15 13:57:54
 * @LastEditTime: 2022-09-22 10:27:39
 * @LastEditors: VSCode
 * @Reference:
 */
package make

import (
	"fmt"
	"os"

	"github.com/laravelGo/core/console"
	"github.com/spf13/cobra"
)

var CmdMakeModel = &cobra.Command{
	Use:   "model",
	Short: "生成模型文件指令，示例: make model user",
	Run:   runMakeModel,
	Args:  cobra.ExactArgs(1), //只有一个参数
}

func runMakeModel(cmd *cobra.Command, args []string) {
	console.Warning("开始生成模型模版文件")
	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(args[0])
	// 确保模型的目录存在，例如 `app/models/user`
	dir := fmt.Sprintf("app/models/%s/", model.PackageName)
	// os.MkdirAll 会确保父目录和子目录都会创建，第二个参数是目录权限，使用 0777
	os.MkdirAll(dir, os.ModePerm)
	// 替换变量
	createFileFromStub(dir+model.PackageName+"_model.go", "model/model", model)
	createFileFromStub(dir+model.PackageName+"_util.go", "model/model_util", model)
	console.Success("模型文件已经生成,请自行修改")
}
