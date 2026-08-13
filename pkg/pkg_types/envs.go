package pkg_types

import (
	"strconv"
	"strings"
	"time"

	"GoRestify/pkg/decimal"
	"GoRestify/pkg/dictionary"
)

// Envkey environment key
type Envkey string

// Envs holds all environments
type Envs map[Envkey]string

// ToBool cast string to boolean
func (p Envs) ToBool(key Envkey) bool {
	return strings.ToUpper(p[key]) == "TRUE"
}

// StrToUint convert string number to RowID
func StrToUint(strNum string) (id uint, err error) {
	tmpID, err := strconv.ParseUint(strNum, 10, 64)
	id = uint(tmpID)
	return
}

// ToUint casting string to uint
func (p Envs) ToUint(key Envkey) uint {
	num, _ := StrToUint(p[key])
	return num
}

// ToPtrUint casting string to uint
func (p Envs) ToPtrUint(key Envkey) *uint {
	num, _ := StrToUint(p[key])
	return &num
}

// ToUint64 casting string to Uint64
func (p Envs) ToUint64(key Envkey) uint64 {
	num, _ := strconv.ParseUint(p[key], 10, 64)
	return num
}

// ToInt64 casting string to Int64
func (p Envs) ToInt64(key Envkey) int64 {
	num, _ := strconv.ParseInt(p[key], 10, 64)
	return num
}

// ToLang casting string to dictionary.Lang
func (p Envs) ToLang(key Envkey) dictionary.Lang {
	lang := dictionary.Lang(p[key])
	return lang
}

// ToByte casting string to []byte
func (p Envs) ToByte(key Envkey) []byte {
	return []byte(p[key])
}

// ToDuration casting string to time.Duration
func (p Envs) ToDuration(key Envkey) time.Duration {
	num := p.ToUint64(key)
	return time.Duration(num)
}

// ToInt casting string to int
func (p Envs) ToInt(key Envkey) int {
	num, _ := strconv.Atoi(p[key])
	return num
}

// ToFloat64 casting string to float64
func (p Envs) ToFloat64(key Envkey) float64 {
	num, _ := strconv.ParseFloat(p[key], 64)
	return num
}

// ToLocation casting string to time.Location
func (p Envs) ToLocation(key Envkey) *time.Location {
	tz, _ := time.LoadLocation(p[key])
	return tz
}

// ToDecimal casting string to decimal
func (p Envs) ToDecimal(key Envkey) *decimal.Decimal {
	num, _ := decimal.NewFromString(p[key])
	return &num
}
