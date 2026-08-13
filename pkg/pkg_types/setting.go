package pkg_types

import (
	"strconv"

	"GoRestify/pkg/decimal"
)

// Setting type
type Setting string

// SettingMap model
type SettingMap struct {
	Value string
}

// ToUint return uint for id
func (p SettingMap) ToUint() uint {
	n, _ := StrToUint(p.Value)
	return n
}

// ToFloat64 return uint for id
func (p SettingMap) ToFloat64() float64 {
	n, _ := strconv.ParseFloat(p.Value, 64)
	return n
}

// ToInt return uint for id
func (p SettingMap) ToInt() int {
	n, _ := strconv.Atoi(p.Value)
	return n
}

// ToBool return bool
func (p SettingMap) ToBool() bool {
	n, _ := strconv.ParseBool(p.Value)
	return n
}

// ToDecimal casting string to decimal
func (p SettingMap) ToDecimal() *decimal.Decimal {
	num, _ := decimal.NewFromString(p.Value)
	return &num
}
