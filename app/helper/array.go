/*
 * @Description:
 * @Author: mali
 * @Date: 2023-07-31 15:10:59
 * @LastEditTime: 2023-07-31 15:14:38
 * @LastEditors: VSCode
 * @Reference:
 */
package helper

func ArrayColumn(input []map[string]interface{}, columnKey string, indexKey string) interface{} {
	if indexKey != "" {
		columns := map[interface{}]interface{}{}
		for _, val := range input {
			if columnKey != "" {
				if v, ok := val[columnKey]; ok {
					columns[val[indexKey]] = v
				}
			} else {
				columns[val[indexKey]] = val
			}
		}
		return columns
	} else {
		columns := make([]interface{}, 0, len(input))
		for _, val := range input {
			if columnKey != "" {
				if v, ok := val[columnKey]; ok {
					columns = append(columns, v)
				}
			} else {
				columns = append(columns, val)
			}
		}
		return columns
	}
}

func InArray(needle interface{}, hystack interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range hystack.([]string) {
			if key == item {
				return true
			}
		}
	case uint:
		for _, item := range hystack.([]uint) {
			if key == item {
				return true
			}
		}
	case uint8:
		for _, item := range hystack.([]uint8) {
			if key == item {
				return true
			}
		}
	case uint16:
		for _, item := range hystack.([]uint16) {
			if key == item {
				return true
			}
		}
	case uint32:
		for _, item := range hystack.([]uint32) {
			if key == item {
				return true
			}
		}
	case uint64:
		for _, item := range hystack.([]uint64) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range hystack.([]int) {
			if key == item {
				return true
			}
		}
	case int8:
		for _, item := range hystack.([]int8) {
			if key == item {
				return true
			}
		}
	case int16:
		for _, item := range hystack.([]int16) {
			if key == item {
				return true
			}
		}
	case int32:
		for _, item := range hystack.([]int32) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range hystack.([]int64) {
			if key == item {
				return true
			}
		}
	case float32:
		for _, item := range hystack.([]float32) {
			if key == item {
				return true
			}
		}
	case float64:
		for _, item := range hystack.([]float64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}

/**
 * @Author: mali
 * @Func:
 * @Description: 获取
 * @Param:
 * @Return:
 * @Example:
 * @param {[]interface{}} arr
 */
func ArrayUnique(arr []interface{}) []interface{} {
	size := len(arr)
	result := make([]interface{}, 0, size)
	temp := map[interface{}]struct{}{}
	for i := 0; i < size; i++ {
		if _, ok := temp[arr[i]]; !ok {
			temp[arr[i]] = struct{}{}
			result = append(result, arr[i])
		}
	}
	return result
}

/**
 * @Author: mali
 * @Func:
 * @Description: 任意类型切片的元素唯一
 * @Param:
 * @Return:
 * @Example:
 * @param {[]interface{}} arr
 */
func ArrayUniqueInt32(input []int32) []int32 {
	size := len(input)
	result := make([]int32, 0, size)
	temp := map[int32]struct{}{}
	for i := 0; i < size; i++ {
		if _, ok := temp[input[i]]; !ok {
			temp[input[i]] = struct{}{}
			result = append(result, input[i])
		}
	}
	return result
}

/**
 * @Author: mali
 * @Func:
 * @Description: 获取索引key
 * @Param:
 * @Return:
 * @Example:
 * @param {map[interface{}]interface{}} elements
 */
func ArrayKeys(elements map[interface{}]interface{}) []interface{} {
	i, keys := 0, make([]interface{}, len(elements))
	for key := range elements {
		keys[i] = key
		i++
	}
	return keys
}
