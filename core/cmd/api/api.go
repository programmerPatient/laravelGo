/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-05 14:28:45
 * @LastEditTime: 2022-11-08 11:15:53
 * @LastEditors: VSCode
 * @Reference:
 */
package api

import (
	"fmt"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/laravelGo/bootstrap"
	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/console"
	"github.com/spf13/cobra"
)

var ApiCmd = &cobra.Command{
	Use:   "api",
	Short: "开启api服务命令",
	Args:  cobra.NoArgs,
	Run:   runApi,
}

/**
 * @Author: mali
 * @Func: runApi
 * @Description: 指令运行主体
 * @Param:
 * @Return:
 * @Example:
 * @param {*cobra.Command} cmd
 * @param {[]string} args
 */
func runApi(cmd *cobra.Command, args []string) {
	api := gin.New()
	bootstrap.Start(api)
	IsDebug := config.GetBool("app.debug")
	addr := config.GetString("app.host") + ":" + config.GetString("app.port")
	if !IsDebug {
		gin.SetMode(gin.ReleaseMode)
	} else {
		console.Success("Listening and serving HTTP on " + addr + "\n")
	}
	//监听端口
	err := endless.ListenAndServe(addr, api)
	if err != nil {
		console.Exit(fmt.Sprintf("启动服务失败: %s\n", err))
	}
}
