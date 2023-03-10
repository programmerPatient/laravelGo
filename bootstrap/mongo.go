/*
 * @Description:
 * @Author: mali
 * @Date: 2023-03-02 14:26:37
 * @LastEditTime: 2023-03-10 14:41:53
 * @LastEditors: VSCode
 * @Reference:
 */
package bootstrap

import (
	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/mongo"
)

func SetMongo() {
	//连接mongodb
	mongo.ConnectMongo(
		config.GetString("mongo.host"),
		config.GetUint64("mongo.port"),
		config.GetString("mongo.username"),
		config.GetString("mongo.password"),
		config.GetString("mongo.db"),
		config.GetUint64("mongo.timeOut"),
		config.GetUint64("mongo.maxNum"),
	)
}
