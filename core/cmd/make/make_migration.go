/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-14 10:43:26
 * @LastEditTime: 2022-09-21 09:51:11
 * @LastEditors: VSCode
 * @Reference:
 */
package make

import (
	"fmt"

	"github.com/laravelGo/core/app"
	"github.com/laravelGo/core/console"

	"github.com/spf13/cobra"
)

var CmdMakeMigration = &cobra.Command{
	Use:   "migration",
	Short: "生成迁移文件, 示例: make migration add_users_table",
	Run:   runMakeMigration,
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func runMakeMigration(cmd *cobra.Command, args []string) {
	// 日期格式化
	timeStr := app.TimenowInTimezone().Format("2006_01_02_150405")
	model := makeModelFromString(args[0])
	fileName := timeStr + "_" + model.PackageName
	filePath := fmt.Sprintf("database/migrations/%s.go", fileName)
	createFileFromStub(filePath, "migration", model, map[string]string{"{{FileName}}": fileName})
	console.Success("迁移文件已经生成，请运行 go mod tidy自动加载包, 使用 `migrate up` 迁移数据库.")
}
