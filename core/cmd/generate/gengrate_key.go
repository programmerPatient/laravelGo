/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-15 17:47:53
 * @LastEditTime: 2022-09-15 18:05:38
 * @LastEditors: VSCode
 * @Reference:
 */
package generate

import (
	"os"
	"strings"

	"github.com/laravelGo/core/console"
	"github.com/laravelGo/core/file"
	"github.com/laravelGo/core/helper"
	"github.com/spf13/cobra"
)

var CmdGenerateKey = &cobra.Command{
	Use:   "key",
	Short: "生成应用的key密钥",
	Run:   RunGenerateKey,
}

func RunGenerateKey(cmd *cobra.Command, args []string) {
	//32位的随机字符串
	randomStr := helper.RandomString(32)
	filePath := ".env"
	// 目标文件已存在
	if !file.Exists(filePath) {
		console.Exit(".env文件不存在了")
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		console.Exit(err.Error())
	}
	// 替换APP_KEY
	stringData := string(data)
	start := strings.Index(stringData, "APP_KEY=")
	if start == -1 {
		console.Exit("不存在APP_KEY配置")
	}
	console.Warning(stringData[start : start+40])
	console.Warning("APP_KEY=" + randomStr)

	newData := strings.ReplaceAll(stringData, stringData[start:start+40], "APP_KEY="+randomStr)
	// 存储到目标文件中
	err = file.Put([]byte(newData), filePath)
	if err != nil {
		console.Exit(err.Error())
	}
}
