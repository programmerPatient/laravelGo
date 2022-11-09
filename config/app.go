/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-07 11:34:30
 * @LastEditTime: 2022-11-09 13:52:30
 * @LastEditors: VSCode
 * @Reference:
 */
package config

import "github.com/laravelGo/core/config"

func init() {
	config.AddConfig("app", func() map[string]interface{} {
		return map[string]interface{}{
			// 应用名称
			"name": config.Env("APP_NAME", "defaultName"),
			// 当前环境，用以区分多环境，一般为 local, stage, production, test
			"env": config.Env("APP_ENV", "production"),
			// 是否进入调试模式
			"debug": config.Env("APP_DEBUG", false),
			// 应用服务host
			"host": config.Env("APP_HOST", "127.0.0.1"),
			// 应用服务端口
			"port": config.Env("APP_PORT", 9000),
			// 加密密钥
			"key": config.Env("APP_KEY", ""),
			// 设置时区
			"timezone": config.Env("TIMEZONE", "Asia/Shanghai"),
			// 是否开启自动生成文档
			"swage_handler": config.Env("SWAGE_HANDLER", false),
		}
	})
}
