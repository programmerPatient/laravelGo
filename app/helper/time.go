/*
 * @Description:
 * @Author: mali
 * @Date: 2023-07-31 15:16:21
 * @LastEditTime: 2023-07-31 15:16:24
 * @LastEditors: VSCode
 * @Reference:
 */
package helper

import (
	"fmt"
	"time"
)

/**
 * @Author: mali
 * @Func:
 * @Description: 秒数转化为分秒字符串
 * @Param:
 * @Return:
 * @Example:
 * @param {int64} se
 */
func SecondToMsString(se int64) string {
	if se > 60 {
		min := se / 60
		second := se % 60
		return fmt.Sprintf("%v分%v秒", min, second)
	} else {
		return fmt.Sprintf("%v秒", se)
	}
}

/**
 * @Author: mali
 * @Func:
 * @Description: 时间字符串转换为时间格式
 * @Param:
 * @Return:
 * @Example:
 * @param {string} format
 * @param {string} input
 */
func StringToTime(format string, input string) (time.Time, error) {
	default_format := []string{
		format,
	}
	if Empty(format) {
		//常规类型
		default_format = []string{
			"2006-01-02 15:04:05 +0800 CST",
			"2006-01-02 15:04:05",
			"2006/01/02 15:04:05",
			"2006-01-2 15:04:05",
			"2006/01/2 15:04:05",
			"2006-1-02 15:04:05",
			"2006/1/02 15:04:05",
			"2006-1-2 15:04:05",
			"2006/1/2 15:04:05",
			"2006-01-02 15:04",
			"2006/01/02 15:04",
			"2006-1-02 15:04",
			"2006/1/02 15:04",
			"2006-01-2 15:04",
			"2006/01/2 15:04",
			"2006-1-2 15:04",
			"2006/1/2 15:04",
			"2006-01-02 15",
			"2006/01/02 15",
			"2006-1-02 15",
			"2006/1/02 15",
			"2006-01-2 15",
			"2006/01/2 15",
			"2006-1-2 15",
			"2006/1/2 15",
			"2006-01-02",
			"2006/01/02",
			"2006-01-2",
			"2006/01/2",
			"2006-1-02",
			"2006/1/02",
			"2006-1-2",
			"2006/1/2",
			"15:04:05",
			"15:04",
			"15",
		}
	}
	var return_time time.Time
	var err error
	for _, v := range default_format {
		return_time, err = time.ParseInLocation(v, input, time.Local)
		if !return_time.IsZero() {
			break
		}
	}
	return return_time, err
}

/**
 * @Author: mali
 * @Func:
 * @Description: 时间格式化
 * @Param:
 * @Return:
 * @Example:
 * @param {time.Time} t
 */
func FormatTime(t time.Time) string {
	now := time.Now()
	diff_time := now.Unix() - t.Unix()
	if diff_time < 60 {
		return fmt.Sprintf("%v秒前", diff_time)
	} else if diff_time < 3600 {
		min := diff_time / 60
		return fmt.Sprintf("%v分钟前", min)
	} else if diff_time < 86400 {
		h := diff_time / 3600
		return fmt.Sprintf("%v小时前", h)
	} else if diff_time < 86400*30 {
		day := diff_time / 86400
		return fmt.Sprintf("%v天前", day)
	} else if diff_time < 86400*30*12 {
		month := diff_time / (86400 * 30)
		return fmt.Sprintf("%v个月前", month)
	} else {
		year := diff_time / (86400 * 30 * 12)
		if (diff_time - year*86400*30*12) >= (86400 * 30) {
			month := (diff_time - year*86400*30*12) / (86400 * 30)
			return fmt.Sprintf("%v年%v个月前", year, month)
		} else {
			return fmt.Sprintf("%v年前", year)
		}
	}
}
