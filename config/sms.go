/*
 * @Description: 短信配置
 * @Author: mali
 * @Date: 2022-10-24 15:22:17
 * @LastEditTime: 2023-08-07 09:58:32
 * @LastEditors: VSCode
 * @Reference:
 */
package config

import "github.com/laravelGo/core/config"

func init() {
	config.AddConfig("sms", func() map[string]interface{} {
		return map[string]interface{}{
			//选取的短信厂商
			"type": "yunpian",
			/***短信厂商配置 start***/
			//云片短信配置
			"yunpian": map[string]interface{}{
				"apikey": "",
				/****短信模版配置****/
				"templete": map[string]interface{}{},
				/****短信模版配置****/
			},
			/***短信厂商配置 end***/
		}
	})
}
