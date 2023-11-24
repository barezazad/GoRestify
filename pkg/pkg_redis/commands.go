package pkg_redis

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"GoRestify/pkg/tx"

	"github.com/go-redis/redis/v8"
)

// RedisCon .
type RedisCon struct {
	client *redis.Client
}

// Set set value by key
func (redisCon RedisCon) Set(key string, value interface{}) (err error) {

	payload := StringFyPayload(value)

	if DisplayLog {
		fmt.Println("Redis SET KEY:", key, ":", compactLog(value))
	}

	err = redisCon.client.Set(Ctx, key, payload, 0).Err()
	return
}

// SetWithTTL set value with TTL
func (redisCon RedisCon) SetWithTTL(key string, value interface{}, ttl int) (err error) {

	payload := StringFyPayload(value)

	if DisplayLog {
		fmt.Println("Redis SET KEY with TTL:", key, ":", compactLog(value))
	}

	err = redisCon.client.Set(Ctx, key, payload, time.Duration(ttl)*time.Second).Err()
	return
}

// Get get by key
func (redisCon RedisCon) Get(key string) string {
	value, err := redisCon.client.Get(Ctx, key).Result()
	if err == redis.Nil {
		err = nil
		return ""
	}

	if DisplayLog {
		fmt.Println("Redis GET KEY:", key, ":", compactLog(value))
	}

	return value
}

// GetCache get by key and parsed response data based on provided model
func (redisCon RedisCon) GetCache(tx tx.Tx, key string, model interface{}) (ok bool) {

	if tx.IsLock {
		return
	}

	value, err := redisCon.client.Get(Ctx, key).Result()
	if err != nil {
		err = nil
		return
	}

	if DisplayLog {
		fmt.Println("Redis GET KEY In Cache:", key, ":", compactLog(value))
	}

	if err = json.Unmarshal([]byte(value), &model); err != nil {
		fmt.Println("Error in parsing json in redis getCache:", err)
		return
	}

	return true
}

// Delete delete by key
func (redisCon RedisCon) Delete(key string) (err error) {
	err = redisCon.client.Del(Ctx, key).Err()
	if err != nil {
		return
	}

	if DisplayLog {
		fmt.Println("Redis DELETE KEY:", key)
	}

	return
}

// DeleteMultiKeys delete via multi keys
// keys = []string{customer-1,customer-2}
func (redisCon RedisCon) DeleteMultiKeys(key ...string) (err error) {
	err = redisCon.client.Del(Ctx, key...).Err()
	if err != nil {
		return
	}

	if DisplayLog {
		fmt.Println("Redis DELETE KEY:", compactLog(key))
	}

	return
}

// KeyExist check to key is exist
// example key = "customer-1"
func (redisCon RedisCon) KeyExist(key string) (exist bool) {
	var existInt int64

	if DisplayLog {
		fmt.Println("Redis KEY Exist:", key)
	}

	existInt, err := redisCon.client.Exists(Ctx, key).Result()
	if err != nil {
		return
	}

	return existInt == 1
}

// KeysPattern to get all keys by pattern
// example key = "customer-*"
func (redisCon RedisCon) KeysPattern(key string) (keys []string) {

	if DisplayLog {
		fmt.Println("Redis KEY Pattern:", key)
	}

	keys, err := redisCon.client.Keys(Ctx, key).Result()
	if err != nil {
		return
	}

	return
}

// ResetCacheByKeyPatten to delete all keys by pattern
// example keyPattern = "customer-*"
func (redisCon RedisCon) ResetCacheByKeyPatten(keyPattern string) {

	keyPattern = fmt.Sprintf("%v*", keyPattern)
	keys := redisCon.KeysPattern(keyPattern)

	j := 0
	for i := 0; i < len(keys); i += 5000 {

		k := j + 5000
		if k > len(keys) {
			k = len(keys)
		}

		keySlice := keys[j:k]
		redisCon.DeleteMultiKeys(keySlice...)
		j = k
	}
}

// FlushDB to flush db, clear all data
func (redisCon RedisCon) FlushDB() (err error) {
	err = redisCon.client.FlushDB(Ctx).Err()
	return
}

// HSet to set data in map
func (redisCon RedisCon) HSet(key string, value interface{}) (int64, error) {

	if DisplayLog {
		fmt.Println("Redis HSet KEY:", key, ":", compactLog(value))
	}

	return redisCon.client.HSet(Ctx, key, value).Result()
}

// HGetAll to fetch data
func (redisCon RedisCon) HGetAll(key string) (res *redis.StringStringMapCmd, err error) {

	if DisplayLog {
		fmt.Println("Redis HGetAll KEY:", key)
	}

	res = redisCon.client.HGetAll(Ctx, key)
	if res.Err() != nil {
		return nil, res.Err()
	}
	return
}

// Redis List

// ListPush push to a list
func (redisCon RedisCon) ListPush(key string, value interface{}) (err error) {

	payload := StringFyPayload(value)

	if DisplayLog {
		fmt.Println("Redis Push KEY:", key, ":", compactLog(value))
	}

	err = redisCon.client.LPush(Ctx, key, payload).Err()

	return
}

// ListGet delete item in a list
func (redisCon RedisCon) ListGet(key string, start, stop int64, model interface{}) (ok bool) {

	result, err := redisCon.client.LRange(Ctx, key, start, stop).Result()
	if err != nil || len(result) == 0 {
		return
	}

	if DisplayLog {
		fmt.Println("Redis ListGet KEY:", key, "start,stop", start, stop, ":", compactLog(result))
	}

	switch model.(type) {

	case []string:
		model = result

	default:
		var data []interface{}
		for _, v := range result {
			var el interface{}
			if err = json.Unmarshal([]byte(v), &el); err != nil {
				fmt.Println("Error in parsing JSON in Redis ListGet:", err)
				return
			}
			data = append(data, el)
		}

		// Convert the []interface{} to the desired type
		dataBytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error in marshaling data to []byte:", err)
			return
		}

		if err = json.Unmarshal(dataBytes, model); err != nil {
			fmt.Println("Error in unmarshal data to model:", err)
			return
		}
	}
	ok = true

	return
}

// ListCount delete item in a list
func (redisCon RedisCon) ListCount(key string) (count uint) {

	listLen, err := redisCon.client.LLen(Ctx, key).Result()
	if err != nil {
		fmt.Println("Error in redis ListCount:", err)
		return
	}
	count = uint(listLen)

	if DisplayLog {
		fmt.Println("Redis ListCount KEY: ", key, "Count: ", count)
	}

	return
}

// GetRedisClient to get redis.client
func (redisCon RedisCon) GetRedisClient() *redis.Client {
	return redisCon.client
}

// StringFyPayload cast payload to string
func StringFyPayload(value interface{}) (dataInString string) {

	switch value.(type) {
	case string:
		dataInString = value.(string)
	case int64:
		dataInString = fmt.Sprint(value.(int64))
	case int:
		dataInString = fmt.Sprint(value.(int))
	case uint:
		dataInString = fmt.Sprint(value.(uint))
	default:
		dataInJSON, err := json.Marshal(value)
		if err != nil {
			dataInJSON = []byte(value.(string))
		}
		dataInString = string(dataInJSON)
	}

	return
}

// compactLog it show just 1000 char in a value, if your value is string or struct
func compactLog(value interface{}) (logValue interface{}) {

	// Use reflection to check the type of the 'value' parameter
	valueType := reflect.TypeOf(value)

	if valueType.Kind() == reflect.String {

		logValueStr := value.(string)
		logValue = logValueStr
		if len(logValueStr) > 1000 {
			logValue = logValueStr[:1000] + " ......"
		}

	} else if valueType.Kind() == reflect.Slice && valueType.Elem().Kind() == reflect.Struct {

		dataInJSON, err := json.Marshal(value)
		if err != nil {
			dataInJSON = []byte(value.(string))
		}
		dataStr := string(dataInJSON)

		logValue = dataStr
		if len(dataStr) > 1000 {
			logValue = dataStr[:1000] + " ......"
		}

	} else if valueType.Kind() == reflect.Slice && valueType.Elem().Kind() == reflect.String {

		array := value.([]string)
		arrayStr := strings.Join(array, ",")

		logValue = arrayStr
		if len(arrayStr) > 1000 {
			logValue = arrayStr[:1000] + " ......"
		}

	} else {
		logValue = value
	}

	return
}
