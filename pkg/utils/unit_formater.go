package utils

import (
	"fmt"
)

// ByteConverter to convert Byte balance to readable balance
func ByteConverter(length int) (out string) {

	// Define size constants
	const (
		TB = 1 << 40 // 1099511627776
		GB = 1 << 30 // 1073741824
		MB = 1 << 20 // 1048576
		KB = 1 << 10 // 1024
	)

	var unit string
	var integerPart int
	var decimalPart float64

	// Determine the unit and compute the integer and decimal parts
	switch {
	case length >= TB:
		unit = "TB"
		integerPart = length / TB
		decimalPart = float64(length%TB) / float64(TB)
	case length >= GB && length > 2684354560:
		unit = "GB"
		integerPart = length / GB
		decimalPart = float64(length%GB) / float64(GB)
	case length >= MB:
		unit = "MB"
		integerPart = length / MB
		decimalPart = float64(length%MB) / float64(MB)
	case length >= KB:
		unit = "KB"
		integerPart = length / KB
		decimalPart = float64(length%KB) / float64(KB)
	default:
		unit = "B"
		integerPart = length
		decimalPart = 0
	}

	// Format the decimal part with two decimal places
	return fmt.Sprintf("%.2f %s", float64(integerPart)+decimalPart, unit)
}
