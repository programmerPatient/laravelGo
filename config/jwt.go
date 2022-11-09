/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-13 14:51:32
 * @LastEditTime: 2022-11-09 13:52:10
 * @LastEditors: VSCode
 * @Reference:
 */
package config

import "github.com/laravelGo/core/config"

func init() {
	config.AddConfig("jwt", func() map[string]interface{} {
		return map[string]interface{}{
			// 密钥
			"key": "",
			// 签名颁发者
			"issuer": "",
			// 签名接收者
			"audience": "",
			// 过期时间，单位是分钟，一般不超过两个小时
			"expire_time": config.Env("JWT_EXPIRE_TIME", 1),
			// 允许刷新时间，单位分钟，86400 为两个月，从 Token 的签名时间算起
			"max_refresh_time": config.Env("JWT_MAX_REFRESH_TIME", 86400),
			// debug 模式下的过期时间，方便本地开发调试
			"debug_expire_time": 86400,
		}
	})
}
