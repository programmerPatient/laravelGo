/*
 * @Description:base64图像处理
 * @Author: mali
 * @Date: 2022-11-07 15:15:27
 * @LastEditTime: 2023-08-07 09:45:54
 * @LastEditors: VSCode
 * @Reference:
 */
package image

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/laravelGo/app/helper"
)

/**
 * @Author: mali
 * @Func:
 * @Description: base64写入文件
 * @Param:
 * @Return:
 * @Example:
 * @param {string} path
 * @param {string} base64_image_content
 */
func Base64WriteFile(ctx context.Context, path string, base64_image_content string) (bool, string) {
	b, _ := regexp.MatchString(`^data:\s*image\/(\w+);base64,`, base64_image_content)
	if !b {
		return false, ""
	}

	re, _ := regexp.Compile(`^data:\s*image\/(\w+);base64,`)
	allData := re.FindAllSubmatch([]byte(base64_image_content), 2)
	fileType := string(allData[0][1]) //图片后缀获取

	base64Str := re.ReplaceAllString(base64_image_content, "")

	date := time.Now().Format("20060102")

	if ok := helper.IsFileExist(path + "/" + date); !ok {
		os.Mkdir(path+"/"+date, 0777)
	}

	curFileStr := strconv.FormatInt(time.Now().UnixNano(), 10)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(99999)
	var file string = path + "/" + date + "/" + helper.Sha1(curFileStr+strconv.Itoa(n)) + "." + fileType
	byte, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return false, ""
	}
	err = ioutil.WriteFile(file, byte, 0777)
	if err != nil {
		return false, ""
	}

	return true, file
}
