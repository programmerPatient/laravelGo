/*
 * @Description:
 * @Author: mali
 * @Date: 2023-07-19 10:20:45
 * @LastEditTime: 2023-08-07 09:57:16
 * @LastEditors: VSCode
 * @Reference:
 */
package config

import "github.com/laravelGo/core/config"

func init() {
	config.AddConfig("pinecone", func() map[string]interface{} {
		return map[string]interface{}{
			"apiKey":     "",   //apiKey
			"namespace":  "",   //命名空间
			"indexName":  "",   //索引
			"projectId":  "",   //项目
			"env":        "",   //环境
			"dimensions": 1536, //向量尺寸

		}
	})
}
