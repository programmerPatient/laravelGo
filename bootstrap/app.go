/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-05 16:26:03
 * @LastEditTime: 2022-09-08 11:58:56
 * @LastEditors: VSCode
 * @Reference:
 */
package bootstrap

import "github.com/gin-gonic/gin"

/**
 * @Author: mali
 * @Func:
 * @Description: web服务需要开启的服务
 * @Param:
 * @Return:
 * @Example:
 * @param {*gin.Engine} e
 */
func Start(e *gin.Engine) {
	//路由服务
	SetupRoute(e)
}
