/*
 * @Description:
 * @Author: mali
 * @Date: 2023-11-28 11:17:51
 * @LastEditTime: 2023-11-28 15:08:56
 * @LastEditors: VSCode
 * @Reference:
 */
package config

import "github.com/laravelGo/core/config"

func init() {
	config.AddConfig("kafka", func() map[string]interface{} {
		return map[string]interface{}{
			//集群地址
			"hosts": []string{"192.168.21.23:9192", "192.168.21.23:9292", "192.168.21.23:9392"},
			"topic": "test",
		}
	})
}
