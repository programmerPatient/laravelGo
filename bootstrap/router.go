/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-05 16:59:05
 * @LastEditTime: 2022-09-20 11:08:55
 * @LastEditors: VSCode
 * @Reference:
 */
package bootstrap

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laravelGo/app/http/middleware"
	"github.com/laravelGo/core/config"
	_ "github.com/laravelGo/docs"
	"github.com/laravelGo/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/**
 * @Author: mali
 * @Func: SetRouter
 * @Description: 路由设置
 * @Param:
 * @Return:
 * @Example:
 * @param {*gin.Engine} e
 */
func SetupRoute(e *gin.Engine) {
	//设置全局中间件
	setGlobalMiddleWare(e)
	//设置路由
	setRouter(e)
	//404路由
	set404Router(e)
}

/**
 * @Author: mali
 * @Func: setGlobalMiddleWare
 * @Description: 全局中间件
 * @Param:
 * @Return:
 * @Example:
 * @param {*gin.Engine} e
 */
func setGlobalMiddleWare(e *gin.Engine) {
	//这里可以添加自定义的全局中间件
	e.Use(
		middleware.Logger(),
		middleware.Recovery(),
	)
}

/**
 * @Author: mali
 * @Func: setRouter
 * @Description: 设置路由
 * @Param:
 * @Return:
 * @Example:
 * @param {*gin.Engine} e
 */
func setRouter(e *gin.Engine) {
	//是否开启文档url
	swagHandler := config.GetBool("app.swage_handler")
	if swagHandler {
		// 文档界面访问URL
		// http://127.0.0.1:8080/swagger/index.html
		e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	//api接口路由
	router.AddApiRouter(e)
}

/**
 * @Author: mali
 * @Func: set404Router
 * @Description: 404请求的处理
 * @Param:
 * @Return:
 * @Example:
 * @param {*gin.Engine} e
 */
func set404Router(e *gin.Engine) {
	e.NoRoute(func(ctx *gin.Context) {
		//默认返回json
		ctx.JSON(http.StatusNotFound, gin.H{
			"error_code":    404,
			"error_message": "路由未定义，请确认 url 和请求方法是否正确。",
		})
	})
}
