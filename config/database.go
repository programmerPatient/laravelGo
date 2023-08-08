/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-08 10:11:07
 * @LastEditTime: 2023-08-08 13:49:18
 * @LastEditors: VSCode
 * @Reference:
 */
package config

import (
	"github.com/laravelGo/core/config"
)

func init() {
	config.AddConfig("database", func() map[string]interface{} {
		return map[string]interface{}{
			//数据库类型 支持 mysql
			"connection": config.Env("DB_CONNECTION", "mysql"),
			//慢日志的记录的起始时间 单位毫秒
			"slow_sql_min_time": config.Env("DB_SLOW_SQL_MIN_EXEC_TIME", 200),

			/**连接池配置***/
			//设置最大空闲连接数
			"max_idle_connections": config.Env("DB_MAX_IDLE_CONNECTIONS", 7),
			//设置最大连接数
			"max_open_connections": config.Env("DB_MAX_OPEN_CONNECTIONS", 10),
			// 设置每个链接的过期时间
			"max_life_seconds": config.Env("DB_MAX_LIFE_SECONDS", 5*60),

			//mysql配置
			"mysql": map[string]interface{}{
				// 数据库连接信息
				"host":     config.Env("MYSQL_DB_HOST", "127.0.0.1"),
				"port":     config.Env("MYSQL_PORT", "3306"),
				"database": config.Env("MYSQL_DATABASE", ""),
				"username": config.Env("MYSQL_USERNAME", ""),
				"password": config.Env("MYSQL_PASSWORD", ""),
				"charset":  "utf8mb4",
			},
		}
	})

	//其他数据库的连接
	config.AddConfig("other_database", func() map[string]interface{} {
		return map[string]interface{}{
			//示例数据库连接
			"gpt_database": map[string]interface{}{
				//数据库类型 支持 mysql
				"connection": "mysql",
				//慢日志的记录的起始时间 单位毫秒
				"slow_sql_min_time": 200,
				/**连接池配置***/
				//设置最大空闲连接数
				"max_idle_connections": 1,
				//设置最大连接数
				"max_open_connections": 2,
				//设置每个链接的过期时间
				"max_life_seconds": 5 * 60,
				//mysql配置
				"mysql": map[string]interface{}{
					// 数据库连接信息
					"host":     "",
					"port":     3306,
					"database": "",
					"username": "",
					"password": "",
					"charset":  "",
				},
			},
		}
	})
}
