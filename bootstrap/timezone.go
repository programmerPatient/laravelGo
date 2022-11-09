/*
 * @Description:
 * @Author: mali
 * @Date: 2022-11-08 11:01:46
 * @LastEditTime: 2022-11-08 11:01:53
 * @LastEditors: VSCode
 * @Reference:
 */
package bootstrap

import (
	"time"

	"github.com/laravelGo/core/config"
)

//全局的时区设置
func SetTimeSone() {
	loc, _ := time.LoadLocation(config.GetString("app.timezone"))
	time.Local = loc
}
