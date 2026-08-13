package utils

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// StrToInt string to int
func StrToInt(value string) (result int, err error) {
	result, err = strconv.Atoi(value)
	return
}

// StrToInt64 string to int64
func StrToInt64(value string) (result int64, err error) {
	result, err = strconv.ParseInt(value, 10, 64)
	return
}

// IntToStr int to string
func IntToStr(value int) (result string) {
	result = strconv.Itoa(value)
	return
}

// UintToStr uint to string
func UintToStr(value uint) (result string) {
	result = strconv.FormatUint(uint64(value), 10)
	return
}

// StrToUint string to uint
func StrToUint(value string) (result uint, err error) {
	tmpID, err := strconv.ParseUint(value, 10, 64)
	result = uint(tmpID)
	return
}

// StrToUint64 string to uint64
func StrToUint64(value string) (result uint64, err error) {
	result, err = strconv.ParseUint(value, 10, 64)
	return
}

// StrToFloat string to float
func StrToFloat(value string) (result float64, err error) {
	result, err = strconv.ParseFloat(value, 64)
	return
}

// IntToPointer int to pointer
func IntToPointer(value int) *int {
	return &value
}

// StrToByte string to byte
func StrToByte(value string) []byte {
	return []byte(value)
}

// ByteToStr byte to string
func ByteToStr(value []byte) string {
	return string(value)
}

// StrToDuration string to duration
func StrToDuration(value string) time.Duration {
	num, _ := StrToUint64(value)
	return time.Duration(num)
}

// StrToBool string to bool
func StrToBool(value string) bool {
	switch strings.ToLower(value) {
	case "0", "false":
		return false
	case "1", "true":
		return true
	default:
		return false
	}
}

// UintPointer returns a *uint
func UintPointer(num uint) *uint {
	return &num
}

// StrPointer returns a *string
func StrPointer(str string) *string {
	return &str
}

// JSONToString to convert json to string
func JSONToString(model interface{}) (result string, err error) {
	encodingData, err := json.Marshal(model)
	result = string(encodingData)
	return
}
