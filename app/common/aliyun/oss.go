package aliyun

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/laravelGo/app/common/image"
	"github.com/laravelGo/app/helper"
	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/logger"
	imgtype "github.com/shamsher31/goimgtype"
)

type AliyunOssClient struct {
	Client     *oss.Client
	Context    context.Context //上下文
	bucket     string
	endpoint   string
	key_id     string
	key_secret string
}

/**
 * @Author: mali
 * @Func:
 * @Description: 只定义配置的oss客户端
 * @Param:
 * @Return:
 * @Example:
 * @param {*} endpoint
 * @param {*} key_id
 * @param {string} key_secret
 */
func GetAliyunOssClient(ctx context.Context, endpoint, key_id, key_secret string) *AliyunOssClient {
	return &AliyunOssClient{
		Client:     getAliyunOssClient(ctx, endpoint, key_id, key_secret),
		key_secret: key_secret,
		key_id:     key_id,
		endpoint:   endpoint,
		Context:    ctx,
	}
}

/**
 * @Author: mali
 * @Func:
 * @Description: 获取oss客户端
 * @Param:
 * @Return:
 * @Example:
 * @param {*} endpoint
 * @param {*} key_id
 * @param {string} key_secret
 */
func getAliyunOssClient(ctx context.Context, endpoint, key_id, key_secret string) *oss.Client {
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, key_id, key_secret)
	if err != nil {
		logger.ErrorJSON("GetOssClient", "获取OSS上传Client报错", err)
		return nil
	}
	return client
}

/**
 * @Author: mali
 * @Func:
 * @Description: 设置当前使用的存储空间
 * @Param:
 * @Return:
 * @Example:
 * @param {string} bucket
 */
func (all_oss *AliyunOssClient) SetBucket(bucket string) {
	all_oss.bucket = bucket
}

/**
 * @Author: mali
 * @Func:
 * @Description: 设置当前使用的域名
 * @Param:
 * @Return:
 * @Example:
 * @param {string} endpoint
 */
func (all_oss *AliyunOssClient) SetEndpoint(endpoint string) {
	all_oss.endpoint = endpoint
}

/**
 * @Author: mali
 * @Func:
 * @Description: 获取指定的存储空间
 * @Param:
 * @Return:
 * @Example:
 */
func (ali_oss *AliyunOssClient) GetOssBucket() *oss.Bucket {
	if helper.Empty(ali_oss.Client) {
		return nil
	}
	bucket := ""
	if helper.Empty(ali_oss.bucket) {
		// yourBucketName填写Bucket名称
		bucket = config.GetString("aliyun.oss.Bucket")
	} else {
		bucket = ali_oss.bucket
	}
	// 指定bucket
	buckets, err := ali_oss.Client.Bucket(bucket)
	if err != nil {
		logger.ErrorJSON("GetOssClient", "获取OSS上传Client报错", err)
		return nil
	}
	return buckets
}

/**
 * @Author: mali
 * @Func:
 * @Description: 文件上传
 * @Param:
 * @Return:
 * @Example:
 * @param {string} path
 * @param {string} file_path
 */
func (ali_oss *AliyunOssClient) UploadFile(path string, file_path string) string {
	bukcket := ali_oss.GetOssBucket()
	fd, err := os.Open(file_path)
	if err != nil {
		return ""
	}
	defer fd.Close()
	err = bukcket.PutObject(path, fd)
	if err != nil {
		logger.ErrorString("aliyun_oss", "aliyun_oss文件上传失败", err.Error())
		return ""
	}
	return fmt.Sprintf("http://%v.%v/%v", ali_oss.bucket, ali_oss.endpoint, path)
}

/**
 * @Author: mali
 * @Func:
 * @Description: 上传base64文件
 * @Param:
 * @Return:
 * @Example:
 * @param {*} base64
 * @param {string} path
 * @param {string} file_path
 */
func (ali_oss *AliyunOssClient) UploadBase64Image(base64, dir, oss_dir string) (string, error) {
	result, file_path := image.Base64WriteFile(ali_oss.Context, dir, base64)
	if result {
		defer func() {
			// 临时文件不用了，需要移除
			err := os.Remove(file_path)
			if err != nil {
				logger.ErrorString("阿里云oss上传base64图片", "base64临时文件删除失败", err.Error())
			}
		}()
		res := strings.Split(file_path, dir)
		oss_path := ali_oss.UploadFile(oss_dir+res[1], file_path)
		if helper.Empty(oss_path) {
			return "", fmt.Errorf("上传失败")
		}
		return oss_path, nil
	}
	return "", fmt.Errorf("上传失败")
}

/**
 * @Author: mali
 * @Func:
 * @Description: 根据url上传图片
 * @Param:
 * @Return:
 * @Example:
 * @param {context.Context} ctx
 * @param {string} url
 */
func (ali_oss *AliyunOssClient) UploadImageByUrl(ctx context.Context, url, upload_base string) string {
	local_image_url, _ := helper.DownLoadImage(url, "storage/tmpfile/")
	datatype, errs := imgtype.Get(local_image_url)
	if errs != nil {
		logger.ErrorString("网络图片类型异常", url, errs.Error())
		return ""
	}
	defer func() {
		//临时文件不用了，需要移除
		err := os.Remove(local_image_url)
		if err != nil {
			logger.ErrorString("网络图片上传的临时文件删除失败", url, err.Error())
		}
	}()
	if strings.Contains(datatype, "image") {
		index := strings.Index("image", datatype)
		datatype = datatype[index+len("image/")+1:]
	}
	upload_base += helper.GetUuid() + "." + datatype
	return ali_oss.UploadFile(upload_base, local_image_url)
}
