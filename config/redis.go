/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-14 10:43:26
 * @LastEditTime: 2022-09-21 14:21:15
 * @LastEditors: VSCode
 * @Reference:
 */
package config

import (
	"github.com/laravelGo/core/config"
)

func init() {

	config.AddConfig("redis", func() map[string]interface{} {
		return map[string]interface{}{

			"host":     config.Env("REDIS_HOST", "127.0.0.1"),
			"port":     config.Env("REDIS_PORT", "6379"),
			"password": config.Env("REDIS_PASSWORD", ""),

			// 业务类存储使用 1 (图片验证码、短信验证码、会话)
			"database": config.Env("REDIS_MAIN_DB", 0),
		}
	})
}
