/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-08 09:37:21
 * @LastEditTime: 2022-09-09 13:49:21
 * @LastEditors: VSCode
 * @Reference:
 */
package bootstrap

import (
	"errors"
	"fmt"
	"time"

	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/database"
	"github.com/laravelGo/core/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDB() {
	var db gorm.Dialector
	connect_type := config.GetString("database.connection")
	switch connect_type {
	case "mysql":
		db = mysql.New(mysql.Config{
			DSN: fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
				config.GetString("database.mysql.username"),
				config.GetString("database.mysql.password"),
				config.GetString("database.mysql.host"),
				config.GetString("database.mysql.port"),
				config.GetString("database.mysql.database"),
				config.GetString("database.mysql.charset"),
			),
		})
	default:
		panic(errors.New("不支持 " + connect_type + "类型的数据库"))
	}
	//使用gorm的默认日志服务
	database.Connect(db, logger.NewGormLogger(time.Duration(config.GetInt64("database.slow_sql_min_time"))))
	// 设置最大连接数
	database.SQLDB.SetMaxOpenConns(config.GetInt("database.max_open_connections"))
	// 设置最大空闲连接数
	database.SQLDB.SetMaxIdleConns(config.GetInt("database.max_idle_connections"))
	// 设置每个链接的过期时间
	database.SQLDB.SetConnMaxLifetime(time.Duration(config.GetInt("database.max_life_seconds")) * time.Second)
}
