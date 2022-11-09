/*
 * @Description:日志开启
 * @Author: mali
 * @Date: 2022-09-08 15:39:46
 * @LastEditTime: 2022-09-08 16:34:33
 * @LastEditors: VSCode
 * @Reference:
 */
package bootstrap

import (
	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/logger"
)

func SetLogger() {
	logger.InitLogger(
		config.GetString("logger.filename"),
		config.GetInt("logger.max_size"),
		config.GetInt("logger.max_backup"),
		config.GetInt("logger.max_age"),
		config.GetBool("logger.compress"),
		config.GetString("logger.type"),
		config.GetString("logger.level"),
	)
}
