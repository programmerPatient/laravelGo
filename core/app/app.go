/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-08 15:28:29
 * @LastEditTime: 2022-09-14 10:43:44
 * @LastEditors: VSCode
 * @Reference:
 */
// Package app 应用信息
package app

import (
	"time"

	"github.com/laravelGo/core/config"
)

func IsLocal() bool {
	return config.GetString("app.env") == "local"
}

func IsProduction() bool {
	return config.GetString("app.env") == "production"
}

func IsTesting() bool {
	return config.GetString("app.env") == "testing"
}

// TimenowInTimezone 获取当前时间，支持时区
func TimenowInTimezone() time.Time {
	chinaTimezone, _ := time.LoadLocation(config.GetString("app.timezone"))
	return time.Now().In(chinaTimezone)
}
