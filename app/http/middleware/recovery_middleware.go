/*
 * @Description: 捕获异常中间价
 * @Author: mali
 * @Date: 2022-09-08 17:05:02
 * @LastEditTime: 2022-11-15 14:30:40
 * @LastEditors: VSCode
 * @Reference:
 */
package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laravelGo/core/logger"
	"go.uber.org/zap"
)

// Recovery 使用 zap.Error() 来记录 Panic 和 call stack
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		//延迟处理的函数
		defer func() {
			// 发生宕机时，获取panic传递的上下文并打印
			if err := recover(); err != nil {
				// 获取用户的请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, true)
				// 链接中断，客户端中断连接为正常行为，不需要记录堆栈信息
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				// 链接中断的情况
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.String("request_id", c.GetString("request_id")), //用于链路追踪使用
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					// 链接已断开，无法写状态码
					return
				}
				// 如果不是链接中断，就开始记录堆栈信息
				logger.Error("recovery from panic",
					zap.String("request_id", c.GetString("request_id")), //用于链路追踪使用
					zap.Time("time", time.Now()),                        //记录时间
					zap.Any("error", err),                               //记录错误信息
					zap.String("request", string(httpRequest)),          //请求信息
					zap.Stack("stacktrace"),                             //调用堆栈信息
				)
				// 返回 500 状态码
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "服务器内部错误，请稍后再试",
				})
			}
		}()
		c.Next()
	}
}
