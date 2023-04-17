/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-05 14:01:03
 * @LastEditTime: 2023-04-17 09:22:23
 * @LastEditors: VSCode
 * @Reference:
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/laravelGo/app/cmd"
	"github.com/laravelGo/bootstrap"
	bsConfig "github.com/laravelGo/config"
	"github.com/laravelGo/core/cmd/api"
	"github.com/laravelGo/core/cmd/cron"
	"github.com/laravelGo/core/cmd/generate"
	"github.com/laravelGo/core/cmd/make"
	"github.com/laravelGo/core/cmd/migrate"
	"github.com/laravelGo/core/cmd/queue"
	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/console"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "根命令",
	// 所有子命令都会执行以下代码
	PersistentPreRun: func(command *cobra.Command, args []string) {
		//配置文件加载
		config.InitConfig()
		//全局时区设置
		bootstrap.SetTimeSone()
		//开启日志
		bootstrap.SetLogger()
		//开启数据库连接
		bootstrap.SetupDB()
		// 初始化 Redis
		bootstrap.SetupRedis()
		//是否开启mongo连接
		if config.GetBool("mongo.status") {
			//初始化mongo
			bootstrap.SetMongo()
		}
	},
}

func init() {
	bsConfig.InitConfig()
}

/**
 * @Author: mali
 * @Func:
 * @Description:
 * @Param:
 * @Return:
 * @Example:
 */
func RunCmd() {
	// 注册子命令
	rootCmd.AddCommand(
		api.ApiCmd,
		migrate.CmdMigrate,
		make.CmdMake,
		generate.CmdGenerate,
		queue.CmdQueue,
		cron.CronCmd,
	)
	//注册自定义子命令
	rootCmd.AddCommand(cmd.InintCmd()...)
	RegisterDefaultCmd(rootCmd, api.ApiCmd)
	if err := rootCmd.Execute(); err != nil {
		//打印错误信息并终止命令
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}

/**
 * @Author: mali
 * @Func: RegisterDefaultCmd
 * @Description: 注册默认命令
 * @Param:
 * @Return:
 * @Example:
 * @param {*cobra.Command} rootCmd
 * @param {*cobra.Command} subCmd 默认运行的指令
 */
func RegisterDefaultCmd(rootCmd *cobra.Command, subCmd *cobra.Command) {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	args := os.Args[1:]
	var firstArg = ""
	if len(args) > 0 {
		firstArg = args[0]
	} else {
		firstArg = ""
	}
	if err == nil && cmd.Use == rootCmd.Use && firstArg != "-h" && firstArg != "--help" {
		args := append([]string{subCmd.Use}, args...)
		rootCmd.SetArgs(args)
	}
}
