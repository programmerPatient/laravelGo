/*
 * @Description:
 * @Author: mali
 * @Date: 2023-07-19 10:11:30
 * @LastEditTime: 2023-08-07 10:07:58
 * @LastEditors: VSCode
 * @Reference:
 */
package pinecone

import (
	"encoding/json"
	"fmt"

	"github.com/laravelGo/app/helper"
	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/logger"
	"github.com/spf13/cast"
)

const (
	PINECONE_VECTORS_UPSERT = "vectors/upsert" //插入更新地址
	PINECONE_QUERY          = "query"
	PINECONE_VECTOR_DELETE  = "vectors/delete" //删除
)

func getUrl(path string) string {
	return fmt.Sprintf(
		"https://%s-%s.svc.%s.pinecone.io/%s",
		config.GetString("pinecone.indexName"),
		config.GetString("pinecone.projectId"),
		config.GetString("pinecone.env"),
		path,
	)
}

func getApiKey() string {
	return config.GetString("pinecone.apiKey")
}
func getNameSpace() string {
	return config.GetString("pinecone.namespace")
}

func VectorUpsert(param map[string]interface{}) (int64, error) {
	if _, ok := param["namespace"]; !ok {
		param["namespace"] = getNameSpace()
	}
	res, err := helper.Http("POST", getUrl(PINECONE_VECTORS_UPSERT), param, map[string]string{
		"accept":       "application/json",
		"content-type": "application/json",
		"Api-Key":      getApiKey(),
	})
	if err != nil {
		logger.ErrorString("pinecone", "VectorUpsert", err.Error())
		return 0, err
	}
	var data map[string]interface{}
	json.Unmarshal([]byte(res), &data)
	if _, ok := data["upsertedCount"]; ok && cast.ToInt64(data["upsertedCount"]) > 0 {
		return cast.ToInt64(data["upsertedCount"]), nil
	} else {
		logger.ErrorString("pinecone", "VectorUpsertError", string(res))
		return 0, fmt.Errorf("向量数据插入到pinecone失败")
	}
}

func VectorQuery(param map[string]interface{}) ([]interface{}, error) {
	if _, ok := param["namespace"]; !ok {
		param["namespace"] = getNameSpace()
	}
	res, err := helper.Http("POST", getUrl(PINECONE_QUERY), param, map[string]string{
		"accept":       "application/json",
		"content-type": "application/json",
		"Api-Key":      getApiKey(),
	})
	if err != nil {
		logger.ErrorString("pinecone", "VectorQuery", err.Error())
		return nil, err
	}
	var data map[string]interface{}
	json.Unmarshal([]byte(res), &data)
	if _, ok := data["matches"]; ok {
		return data["matches"].([]interface{}), nil
	} else {
		logger.ErrorString("pinecone", "VectorQueryError", string(res))
		return nil, fmt.Errorf("向量数据pinecone查询失败")
	}
}

/**
 * @Author: mali
 * @Func:
 * @Description: 根据id删除向量数据
 * @Param:
 * @Return:
 * @Example:
 * @param {[]string} id
 */
func VectorDelete(param map[string]interface{}) (bool, error) {
	if _, ok := param["namespace"]; !ok {
		param["namespace"] = getNameSpace()
	}
	res, err := helper.Http("POST", getUrl(PINECONE_VECTOR_DELETE), param, map[string]string{
		"accept":       "application/json",
		"content-type": "application/json",
		"Api-Key":      getApiKey(),
	})
	if err != nil {
		logger.ErrorString("pinecone", "VectorDeleteError", err.Error())
		return false, err
	}
	if res == "{}" {
		return true, nil
	} else {
		logger.ErrorString("pinecone", "VectorDeleteError", string(res))
		return false, fmt.Errorf("向量数据pinecone删除	失败")
	}
}

/**
 * @Author: mali
 * @Func:
 * @Description: 数据迁移
 * @Param:
 * @Return:
 * @Example:
 * @param {*} from_namespace 从一个namespace迁移到另一个namespace
 * @param {string} to_namespace
 * @param {int64} topK
 */
func Migrate(from_namespace, to_namespace string, topK int64) {
	vectors := []map[string]interface{}{}
	vector := []float32{}
	for i := 0; i < config.GetInt("pinecone.dimensions"); i++ {
		vector = append(vector, 0.1)
	}
	data, _ := VectorQuery(map[string]interface{}{
		"topK":            topK,
		"includeMetadata": true,
		"includeValues":   true,
		"vector":          vector,
		"namespace":       from_namespace,
	})
	id := []interface{}{}
	for _, v := range data {
		data := v.(map[string]interface{})
		metadata := data["metadata"].(map[string]interface{})
		id = append(id, metadata["id"])
		if len(vectors) >= 100 {
			param := map[string]interface{}{
				"namespace": to_namespace,
				"vectors":   vectors,
			}
			_, err := VectorUpsert(param)
			if err != nil {
				fmt.Printf("迁移失败的id为%+v，\n错误信息为 %v", id, err.Error())
			} else {
				fmt.Printf("成功的迁移的id为%+v", id)
			}
			id = []interface{}{}
			vectors = []map[string]interface{}{}
		}
		vectors = append(vectors, data)
	}
	if !helper.Empty(vectors) {
		param := map[string]interface{}{
			"namespace": to_namespace,
			"vectors":   vectors,
		}
		_, err := VectorUpsert(param)
		if err != nil {
			fmt.Printf("迁移失败的id为%+v，\n错误信息为 %v", id, err.Error())
		} else {
			fmt.Printf("成功的迁移的id为%+v", id)
		}
	}
}
