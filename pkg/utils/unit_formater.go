package utils

import (
	"fmt"
	"strconv"
)

// ByteConverter to convert Byte balance to readable balance
func ByteConverter(length int) (out string) {

	// sizes
	const (
		TB = 1099511627776
		GB = 1073741824
		MB = 1048576
		KB = 1024
	)

	var unit string
	var i int
	var remainder int
	decimals := 2

	// Get whole number, and the remainder for decimals
	switch {
	case length > TB:
		unit = "TB"
		i = length / TB
		remainder = length - (i * TB)
		break
	case length > GB:
		unit = "GB"
		i = length / GB
		remainder = length - (i * GB)
		break
	case length > MB:
		unit = "MB"
		i = length / MB
		remainder = length - (i * MB)
		break
	default:
		unit = "KB"
		i = length / KB
		remainder = length - (i * KB)
		break
	}

	// This is to calculate missing leading zeroes
	width := 0
	if remainder > GB {
		width = 12
	} else if remainder > MB {
		width = 9
	} else {
		width = 6
	}

	// Insert missing leading zeroes
	remainderString := strconv.Itoa(remainder)
	for iter := len(remainderString); iter < width; iter++ {
		remainderString = "0" + remainderString
	}
	if decimals > len(remainderString) {
		decimals = len(remainderString)
	}

	return fmt.Sprintf("%d.%s %s", i, remainderString[:decimals], unit)
}
