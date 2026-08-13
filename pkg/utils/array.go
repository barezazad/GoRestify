package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ArrayIncludes used for checking is item exist in the array or not
func ArrayIncludes(arrayType interface{}, item interface{}) (result bool, err error) {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		err = errors.New("invalid data-type")
		return
	}

	for i := 0; i < arr.Len(); i++ {
		if fmt.Sprintf("%v", arr.Index(i).Interface()) == fmt.Sprintf("%v", item) {
			return true, nil
		}
	}

	return
}

// ArrayDifference between two array
func ArrayDifference(array1 []uint, array2 []uint) []uint {
	var diff []uint
	for i := 0; i < 2; i++ {
		for _, s1 := range array1 {
			found := false
			for _, s2 := range array2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// not found. We add it to return array
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the arrays, only if it was the first loop
		if i == 0 {
			array1, array2 = array2, array1
		}
	}
	return diff
}

// ArrayRemoveDuplicate to remove duplicate in array
func ArrayRemoveDuplicate(array []string) []string {

	inResult := make(map[string]bool)
	var result []string

	for _, el := range array {
		if _, ok := inResult[el]; !ok {
			inResult[el] = true
			el = strings.TrimSpace(el)
			if el != "" {
				result = append(result, el)
			}
		}
	}
	return result
}

// ArrayRemoveElement to remove element in array
func ArrayRemoveElement(array []string, el string) []string {

	for i, v := range array {
		if v == el {
			return append(array[:i], array[i+1:]...)
		}
	}
	return array
}

// ArrayRemoveDuplicateUint to remove duplicate in array
func ArrayRemoveDuplicateUint(array []uint) []string {

	inResult := make(map[uint]bool)
	var result []string

	for _, el := range array {
		if _, ok := inResult[el]; !ok {
			inResult[el] = true
			result = append(result, fmt.Sprint(el))
		}
	}
	return result
}

// ArrayRemoveElementUint to remove element in array
func ArrayRemoveElementUint(array []uint, el uint) []uint {

	for i, v := range array {
		if v == el {
			return append(array[:i], array[i+1:]...)
		}
	}
	return array
}
