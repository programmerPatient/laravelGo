/*
 * @Description:
 * @Author: mali
 * @Date: 2023-06-19 14:25:08
 * @LastEditTime: 2023-08-07 09:56:28
 * @LastEditors: VSCode
 * @Reference:
 */
package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/laravelGo/app/helper"
	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/logger"
	"github.com/spf13/cast"
)

const (
	EMBEDDINGS  = "https://api.openai.com/v1/embeddings"
	COMPLETIONS = "https://api.openai.com/v1/completions"
)

func getApiToken() string {
	return fmt.Sprintf("Bearer %v", config.GetString("openai.apiKey"))
}

func getProxy() string {
	return config.GetString("openai.proxyUrl")
}

//文本转化为向量
func Embeddings(ctx context.Context, input string) (map[string]interface{}, error) {
	data, err := helper.NewProxyHttp("POST", EMBEDDINGS, getProxy(), nil, map[string]interface{}{
		"input": input,
		"model": "text-embedding-ada-002",
	}, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": getApiToken(),
	})
	if err != nil {
		logger.ErrorString("openai", "EmbeddingsError", err.Error())
		return nil, err
	}
	var res map[string]interface{}
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		return nil, err
	}
	if _, ok := res["error"]; ok {
		logger.ErrorJSON("openai", "EmbeddingsError", res)
		return nil, fmt.Errorf(cast.ToString(res["message"]))
	}
	return res, err
}

// 定义请求的结构体
type completionRequest struct {
	Prompt           string  `json:"prompt"`            // prompt 是包含用户问题和数据文本的字符串
	MaxTokens        int     `json:"max_tokens"`        // max_tokens 是生成的文本的最大长度
	Temperature      float64 `json:"temperature"`       // temperature 是控制生成文本的随机性的参数
	N                int     `json:"n"`                 // n 是生成的文本的数量
	Stop             string  `json:"stop"`              // stop 是生成文本的终止条件
	FrequencyPenalty float64 `json:"frequency_penalty"` // frequency_penalty 是控制生成文本中重复词汇的程度的参数
	PresencePenalty  float64 `json:"presence_penalty"`  // presence_penalty 是控制生成文本中出现频率低的词汇的程度的参数
	Model            string  `json:"model"`             // model 是要使用的 GPT 模型的名称
}

// 定义响应的结构体
type completionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

func Completions(ctx context.Context, prompt string) (string, error) {
	data, err := helper.NewProxyHttp("POST", COMPLETIONS, getProxy(), nil, completionRequest{
		Prompt:           prompt,
		MaxTokens:        1024,
		Temperature:      0.5,
		N:                1,
		Stop:             "",
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Model:            "gpt-3.5-turbo",
	}, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": getApiToken(),
	})
	if err != nil {
		logger.ErrorString("openai", "Completions", err.Error())
		return "", err
	}
	var response completionResponse
	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		return "", err
	}
	// 返回生成的答案
	return strings.TrimSpace(response.Choices[0].Text), nil
}
