/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-08 09:37:21
 * @LastEditTime: 2023-08-08 13:50:25
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

/**
 * @Author: mali
 * @Func:
 * @Description: 连接其他的数据库
 * @Param:
 * @Return:
 * @Example:
 * @param {string} database_type
 */
func SetupOtherDB(database_type string) {
	var db gorm.Dialector
	connect_type := config.GetString("other_database." + database_type + ".connection")
	switch connect_type {
	case "mysql":
		db = mysql.New(mysql.Config{
			DSN: fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
				config.GetString("other_database."+database_type+".mysql.username"),
				config.GetString("other_database."+database_type+".mysql.password"),
				config.GetString("other_database."+database_type+".mysql.host"),
				config.GetString("other_database."+database_type+".mysql.port"),
				config.GetString("other_database."+database_type+".mysql.database"),
				config.GetString("other_database."+database_type+".mysql.charset"),
			),
		})
	default:
		panic(errors.New("不支持 " + connect_type + "类型的数据库"))
	}
	//使用gorm的默认日志服务
	database.OtherConnect(database_type, db, logger.NewGormLogger(time.Duration(config.GetInt64("other_database."+database_type+".slow_sql_min_time"))))
	// 设置最大连接数
	database.GetSQLDB(database_type).SetMaxOpenConns(config.GetInt("other_database." + database_type + ".max_open_connections"))
	// 设置最大空闲连接数
	database.GetSQLDB(database_type).SetMaxIdleConns(config.GetInt("other_database." + database_type + ".max_idle_connections"))
	// 设置每个链接的过期时间
	database.GetSQLDB(database_type).SetConnMaxLifetime(time.Duration(config.GetInt("other_database."+database_type+".max_life_seconds")) * time.Second)
}
