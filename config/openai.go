/*
 * @Description:
 * @Author: mali
 * @Date: 2023-08-07 09:56:50
 * @LastEditTime: 2023-08-07 09:57:01
 * @LastEditors: VSCode
 * @Reference:
 */
package config

import "github.com/laravelGo/core/config"

func init() {
	config.AddConfig("openai", func() map[string]interface{} {
		return map[string]interface{}{
			"apiKey":   "", //openai的key
			"proxyUrl": "", //代理地址
		}
	})
}
