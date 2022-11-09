/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-14 10:43:26
 * @LastEditTime: 2022-09-15 14:41:53
 * @LastEditors: VSCode
 * @Reference:
 */
package migrate

import (
	"github.com/laravelGo/core/migrate"
	"github.com/laravelGo/database/migrations"

	"github.com/spf13/cobra"
)

func init() {
	CmdMigrate.AddCommand(
		CmdMigrateUp,
		CmdMigrateRollback,
		CmdMigrateRefresh,
		CmdMigrateReset,
		CmdMigrateFresh,
	)
}

var CmdMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "运行迁移相关的指令",
	// 所有 migrate 下的子命令都会执行以下代码
}

var CmdMigrateUp = &cobra.Command{
	Use:   "up",
	Short: "执行迁移，示例: migrate up",
	Run:   runUp,
}

var CmdMigrateRollback = &cobra.Command{
	Use: "down",
	// 设置别名 migrate down == migrate rollback
	Aliases: []string{"rollback"},
	Short:   "执行回滚 gun，示例: migrate down",
	Run:     runDown,
}

var CmdMigrateReset = &cobra.Command{
	Use:   "reset",
	Short: "Rollback all database migrations",
	Run:   runReset,
}

var CmdMigrateRefresh = &cobra.Command{
	Use:   "refresh",
	Short: "Reset and re-run all migrations",
	Run:   runRefresh,
}

var CmdMigrateFresh = &cobra.Command{
	Use:   "fresh",
	Short: "Drop all tables and re-run all migrations",
	Run:   runFresh,
}

func migrator() *migrate.Migrator {
	// 注册 database/migrations 下的所有迁移文件
	migrations.Initialize()
	// 初始化 migrator
	return migrate.NewMigrator()
}

func runUp(cmd *cobra.Command, args []string) {
	migrator().Up()
}

func runDown(cmd *cobra.Command, args []string) {
	migrator().Rollback()
}

func runReset(cmd *cobra.Command, args []string) {
	migrator().Reset()
}

func runRefresh(cmd *cobra.Command, args []string) {
	migrator().Refresh()
}

func runFresh(cmd *cobra.Command, args []string) {
	migrator().Fresh()
}
