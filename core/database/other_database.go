/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-07 16:53:04
 * @LastEditTime: 2023-08-08 13:48:07
 * @LastEditors: VSCode
 * @Reference:
 */
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/laravelGo/core/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var otherDBSMap sync.Map
var otherSQLDBSMap sync.Map

func GetDB(database_type string) *gorm.DB {
	db, ok := otherDBSMap.Load(database_type)
	if ok {
		return db.(*gorm.DB)
	} else {
		return nil
	}
}

func GetSQLDB(database_type string) *sql.DB {
	db, ok := otherSQLDBSMap.Load(database_type)
	if ok {
		return db.(*sql.DB)
	} else {
		return nil
	}
}

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
func OtherConnect(database_type string, db gorm.Dialector, logger logger.Interface) {
	//连接指定的数据库
	OtherDB, err := gorm.Open(db, &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// 获取通用数据库对象 sql.DB，然后使用其提供的功能
	OtherSQLDB, err := OtherDB.DB()
	otherDBSMap.Store(database_type, OtherDB)
	otherSQLDBSMap.Store(database_type, OtherSQLDB)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func OtherCurrentDatabase(database_type string) (dbname string) {
	OtherDB := GetDB(database_type)
	dbname = OtherDB.Migrator().CurrentDatabase()
	return
}

func OtherDeleteAllTables(database_type string) error {
	var err error
	switch config.GetString("other_database." + database_type + ".connection") {
	case "mysql":
		err = OtherDeleteMySQLTables(database_type)
	case "sqlite":
		err := OtherDeleteAllSqliteTables(database_type)
		if err != nil {
			return err
		}
	default:
		panic(errors.New("other_database connection not supported"))
	}

	return err
}

func OtherDeleteAllSqliteTables(database_type string) error {
	tables := []string{}
	OtherDB := GetDB(database_type)
	// 读取所有数据表
	err := OtherDB.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'").Error
	if err != nil {
		return err
	}

	// 删除所有表
	for _, table := range tables {
		err := OtherDB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func OtherDeleteMySQLTables(database_type string) error {
	dbname := OtherCurrentDatabase(database_type)
	tables := []string{}
	OtherDB := GetDB(database_type)
	// 读取所有数据表
	err := OtherDB.Table("information_schema.tables").
		Where("table_schema = ?", dbname).
		Pluck("table_name", &tables).
		Error
	if err != nil {
		return err
	}

	// 暂时关闭外键检测
	OtherDB.Exec("SET foreign_key_checks = 0;")

	// 删除所有表
	for _, table := range tables {
		err := OtherDB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}

	// 开启 MySQL 外键检测
	OtherDB.Exec("SET foreign_key_checks = 1;")
	return nil
}
