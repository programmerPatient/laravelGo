/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-07 16:53:04
 * @LastEditTime: 2022-09-08 11:21:13
 * @LastEditors: VSCode
 * @Reference:
 */
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/laravelGo/core/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var SQLDB *sql.DB

/**
 * @Author: mali
 * @Func:
 * @Description:
 * @Param:
 * @Return:
 * @Example:
 * @param {gorm.Dialector} db 要链接的数据库
 * @param {logger.Interface} logger 日志
 */
func Connect(db gorm.Dialector, logger logger.Interface) {
	//连接指定的数据库
	var err error
	DB, err = gorm.Open(db, &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// 获取通用数据库对象 sql.DB，然后使用其提供的功能
	SQLDB, err = DB.DB()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func CurrentDatabase() (dbname string) {
	dbname = DB.Migrator().CurrentDatabase()
	return
}

func DeleteAllTables() error {
	var err error
	switch config.GetString("database.connection") {
	case "mysql":
		err = deleteMySQLTables()
	case "sqlite":
		err := deleteAllSqliteTables()
		if err != nil {
			return err
		}
	default:
		panic(errors.New("database connection not supported"))
	}

	return err
}

func deleteAllSqliteTables() error {
	tables := []string{}

	// 读取所有数据表
	err := DB.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'").Error
	if err != nil {
		return err
	}

	// 删除所有表
	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteMySQLTables() error {
	dbname := CurrentDatabase()
	tables := []string{}

	// 读取所有数据表
	err := DB.Table("information_schema.tables").
		Where("table_schema = ?", dbname).
		Pluck("table_name", &tables).
		Error
	if err != nil {
		return err
	}

	// 暂时关闭外键检测
	DB.Exec("SET foreign_key_checks = 0;")

	// 删除所有表
	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}

	// 开启 MySQL 外键检测
	DB.Exec("SET foreign_key_checks = 1;")
	return nil
}
