package tool

import (
	"encoding/json"
	"math/rand"
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
	// 生成一个 300 到 599 秒之间的随机时间
	return time.Duration(300+rand.Intn(300)) * time.Second
}
