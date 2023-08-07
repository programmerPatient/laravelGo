package helper

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/spf13/cast"
)

func HttpWithCookies(method string, url string, data interface{}, cookies []*http.Cookie, headers ...map[string]string) (string, error) {
	client := &http.Client{}
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Length", cast.ToString(req.ContentLength))
	for _, v := range cookies {
		req.AddCookie(v)
	}
	if len(headers) > 0 {
		for key, header := range headers[0] {
			req.Header.Set(key, header)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//关键的一步，去除最后的空格防止读取报错invalid header field name
	body_string := strings.TrimSpace(string(body))
	return body_string, nil
}

func Http(method string, url string, data interface{}, headers ...map[string]string) (string, error) {
	client := &http.Client{}
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Length", cast.ToString(req.ContentLength))
	if len(headers) > 0 {
		for key, header := range headers[0] {
			req.Header.Set(key, header)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//关键的一步，去除最后的空格防止读取报错invalid header field name
	body_string := strings.TrimSpace(string(body))
	return body_string, nil
}

/**
 * @Author: mali
 * @Func:
 * @Description: 代理请求
 * @Param:
 * @Return:
 * @Example:
 * @param {string} method
 * @param {string} urls
 * @param {string} proxy
 * @param {interface{}} data
 * @param {...map[string]string} headers
 */
func ProxyHttp(method string, urls string, proxy string, data interface{}, headers ...map[string]string) (string, error) {
	// 设置代理IP和端口
	proxyUrl, err := url.Parse("https://" + proxy)
	if err != nil {
		return "", err
	}
	proxy_req := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy_req.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		// 处理代理服务器错误
		fmt.Printf(err.Error())
	}

	jsonStr, _ := json.Marshal(data)
	// 创建HTTP客户端并设置代理
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyUrl),
		},
	}
	// 创建HTTP请求
	request, err := http.NewRequest(method, urls, bytes.NewReader(jsonStr))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Length", cast.ToString(request.ContentLength))
	if len(headers) > 0 {
		for key, header := range headers[0] {
			request.Header.Set(key, header)
		}
	}
	// 发送HTTP GET请求并使用代理IP
	response, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}

	// 设置自定义请求头

	// 读取响应内容
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	//关键的一步，去除最后的空格防止读取报错invalid header field name
	body_string := strings.TrimSpace(string(body))
	return body_string, nil
}

/**
 * @Author: mali
 * @Func:
 * @Description: 代理请求
 * @Param:
 * @Return:
 * @Example:
 * @param {string} method
 * @param {string} urls
 * @param {string} proxy
 * @param {interface{}} data
 * @param {...map[string]string} headers
 */
func NewProxyHttp(method string, urls string, proxy string, cookies []*http.Cookie, data interface{}, headers ...map[string]string) (string, error) {
	jsonStr, _ := json.Marshal(data)
	payload := bytes.NewReader(jsonStr)
	client := &http.Client{}
	if !Empty(proxy) {
		// 设置代理IP和端口
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return "", err
		}
		proxy_req := httputil.NewSingleHostReverseProxy(proxyUrl)
		proxy_req.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			// 处理代理服务器错误
			fmt.Printf(err.Error())
		}

		// 创建HTTP客户端并设置代理
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				Proxy:           http.ProxyURL(proxyUrl),
			},
		}
	}

	req, err := http.NewRequest(method, urls, payload)

	if err != nil {
		return "", err
	}
	for _, v := range cookies {
		req.AddCookie(v)
	}
	if len(headers) > 0 {
		for key, header := range headers[0] {
			req.Header.Add(key, header)
		}
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	//关键的一步，去除最后的空格防止读取报错invalid header field name
	body_string := strings.TrimSpace(string(body))
	return body_string, nil
}

//这是HttpsGet请求
func HttpsGet(url string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}} //如果需要测试自签名的证书 这里需要设置跳过证书检测 否则编译报错
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func HttpsPostForm(url string, data url.Values) (string, error) {
	resp, err := http.PostForm(url, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

/**
 * @Author: mali
 * @Func:
 * @Description: get url参数组装
 * @Param:
 * @Return:
 * @Example:
 * @param {map[string]string} params
 */
func HttpBuildQuery(params map[string]string) (param_str string) {
	params_arr := make([]string, 0, len(params))
	for k, v := range params {
		params_arr = append(params_arr, fmt.Sprintf("%s=%s", k, v))
	}
	param_str = strings.Join(params_arr, "&")
	return param_str
}
