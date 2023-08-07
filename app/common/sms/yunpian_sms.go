/*
 * @Description:云片短信发送
 * @Author: mali
 * @Date: 2022-10-21 17:01:27
 * @LastEditTime: 2023-08-07 09:58:01
 * @LastEditors: VSCode
 * @Reference:
 */
package sms

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/laravelGo/app/helper"
	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/logger"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

const url_single_send_sms = "https://sms.yunpian.com/v2/sms/single_send.json"

type YunPianSMS struct {
	ApiKey string
}

func GetYunPianSms() *YunPianSMS {
	return &YunPianSMS{
		ApiKey: config.GetString("sms.yunpian.apikey"),
	}
}

func (sm *YunPianSMS) ShendMesssage(ctx context.Context, mobile string, message string) bool {
	data_send_sms := url.Values{"apikey": {sm.ApiKey}, "mobile": {mobile}, "text": {message}}
	res, err := helper.HttpsPostForm(url_single_send_sms, data_send_sms)
	if err != nil {
		logger.ErrorString("YunPianSMS", "短信发送失败", err.Error())
		return false
	}
	data := map[string]interface{}{}
	err = json.Unmarshal([]byte(res), &data)
	if err == nil {
		if v, ok := data["code"]; ok && cast.ToInt(v) == 0 {
			return true
		} else {
			logger.Warn("YunPianSMS_err_msg", zap.Any("result", res))
		}
	}
	return false
}
