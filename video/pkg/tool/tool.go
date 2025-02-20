package tool

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
)

// 将结构体转换为 map[string]interface{}
func StructToMap(obj interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// redis随机过期时间
func GetRandomExpireTime() time.Duration {
	return time.Duration(300+rand.Intn(300)) * time.Second
}

func StrSliceToInt64Slice(ids []string) []int64 {

	var result []int64
	for _, id := range ids {
		i, _ := strconv.ParseInt(id, 10, 64)
		result = append(result, i)
	}

	return result
}

func Int64SliceToStrSlice(ids []int64) []string {
	var result []string
	for _, id := range ids {
		result = append(result, strconv.FormatInt(id, 10))
	}
	return result
}
