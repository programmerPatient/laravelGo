/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-05 17:05:13
 * @LastEditTime: 2022-11-08 10:57:38
 * @LastEditors: VSCode
 * @Reference:
 */
package controller

import (
	"github.com/gin-gonic/gin"
)

type IndexController struct{}

var IndexC = &IndexController{}

// @Summary 首页接口
// @Description 测试
// @BasePath /api
// @Tags 用户信息
// @Success 200 {string} json "{"message":"success"}"
// @Router / [get]
func (c *IndexController) Index(ctx *gin.Context) {

}
