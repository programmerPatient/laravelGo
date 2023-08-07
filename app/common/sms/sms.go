/*
 * @Description:短信操作类
 * @Author: mali
 * @Date: 2022-10-24 15:12:43
 * @LastEditTime: 2023-08-07 09:57:55
 * @LastEditors: VSCode
 * @Reference:
 */
package sms

import (
	"context"
	"fmt"
	"sync"

	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/logger"
)

//短信接口体接口
type ShortMessage interface {
	ShendMesssage(ctx context.Context, mobile, message string) bool //单条短信发送
}

//短信的结构体
type Message struct {
	Template string
	Data     map[string]string
	Content  string
}

// SMS 是我们发送短信的操作类
type SMS struct {
	ShortMessage
}

//代码运行中需要的时候执行，且只执行一次,保证线程安全
var once sync.Once

var internalSMS *SMS

func GetSMS() *SMS {
	once.Do(func() {
		sms_type := config.GetString("sms.type")
		var sms ShortMessage
		switch sms_type {
		case "yunpian":
			sms = GetYunPianSms()
		default:
			logger.Error(fmt.Sprintf("非法的短信厂商类型配置：%v", sms_type))
		}
		internalSMS = &SMS{
			ShortMessage: sms,
		}
	})
	return internalSMS
}

/**
 * @Author: mali
 * @Func:
 * @Description: 单条短信发送
 * @Param:
 * @Return:
 * @Example:
 * @param {*} mobile
 * @param {string} message
 */
func (s *SMS) ShendMesssage(ctx context.Context, mobile, message string) bool {
	return s.ShortMessage.ShendMesssage(ctx, mobile, message)
}
