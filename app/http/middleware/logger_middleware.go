/*
 * @Description:gin全局的请求中间件
 * @Author: mali
 * @Date: 2022-09-08 17:00:13
 * @LastEditTime: 2022-11-08 11:03:37
 * @LastEditors: VSCode
 * @Reference:
 */
package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laravelGo/app/helper"
	"github.com/laravelGo/core/logger"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 response 内容
		w := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w
		// 获取请求数据
		var requestBody []byte
		if c.Request.Body != nil {
			// c.Request.Body 是一个 buffer 对象，只能读取一次
			requestBody, _ = ioutil.ReadAll(c.Request.Body)
			// 读取后，重新赋值 c.Request.Body ，以供后续的其他操作
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}
		// 设置开始时间
		start := time.Now()
		//执行后续操作
		c.Next()

		/**执行完毕 开始记录日志的逻辑**/
		//执行时间
		cost := time.Since(start)
		//返回状态
		responStatus := c.Writer.Status()
		//日志格式
		logFields := []zap.Field{
			zap.Int("status", responStatus),                                    //状态
			zap.String("time", helper.MicrosecondsStr(cost)),                   //执行时间
			zap.String("request", c.Request.Method+" "+c.Request.URL.String()), //请求的method和url
			zap.String("query", c.Request.URL.RawQuery),                        //请求的参数
			zap.String("ip", c.ClientIP()),                                     //客户端ip
			zap.String("request body", string(requestBody)),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), //gin报错日志详情
		}
		//日志添加头部信息
		for k, v := range c.Request.Header {
			logFields = append(logFields, zap.String(k, v[0]))
		}
		//对数据有影响的操作记录详细日志
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			// 响应的内容
			logFields = append(logFields, zap.String("Response Body", w.body.String()))
		}

		if responStatus > 400 && responStatus <= 499 {
			// 除了 StatusBadRequest 以外，warning 提示一下，常见的有 403 404，开发时都要注意
			logger.Warn("HTTP Warning "+cast.ToString(responStatus), logFields...)
		} else if responStatus >= 500 && responStatus <= 599 {
			// 除了内部错误，记录 error
			logger.Error("HTTP Error "+cast.ToString(responStatus), logFields...)
		} else {
			//成功日志
			logger.Debug("HTTP Access Log", logFields...)
		}

	}
}
