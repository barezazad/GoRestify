package pkg_sql

import (
	"GoRestify/pkg/pkg_err"
	"fmt"
	"strings"
)

// Parser will break the filter into the sub-query
func Parser(str string, cols []string) (string, error) {

	column := make(map[string]string)

	result := strings.ReplaceAll(str, " ", "_SPACE_")
	result = strings.ReplaceAll(result, "\"", "'")
	result = strings.ReplaceAll(result, "(", "PARENTHESES_OPEN ")
	result = strings.ReplaceAll(result, ")", " PARENTHESES_CLOSE")

	for k, v := range filterKeys {
		result = strings.ReplaceAll(result, k, v)
	}
	arr := strings.Split(result, " ")

	for _, v := range cols {

		if strings.Contains(v, ".") {
			splitString := strings.Split(v, ".")
			column[splitString[1]] = v
		}

		if strings.Contains(v, " as ") || strings.Contains(v, " AS ") {
			splitString := strings.Split(v, " ")
			column[splitString[2]] = splitString[0]
		}
	}

	if len(arr) == 0 {
		return "", fmt.Errorf("filter is not valid")
	}

	var col, operator string

	for i, v := range arr {
		v = strings.TrimSpace(v)
		v = strings.ReplaceAll(v, "_SPACE_", " ")
		v = strings.ReplaceAll(v, "PARENTHESES_OPEN", "( ")
		v = strings.ReplaceAll(v, "PARENTHESES_CLOSE", " )")

		arr[i] = v

		if v == "( " || v == " )" {
			continue
		}

		if reverseFilterKeys[v] && col == "" {
			continue
		}

		switch {

		case col == "":
			col = v
			arr[i] = column[v]

			if column[v] == "" {
				err := pkg_err.AddInvalidParam(fmt.Errorf("column '%s' not valid", col), col,
					"column %v not not valid", col)
				return "", err
			}

		case operator == "" && col != "":
			if !reverseFilterKeys[v] {
				err := pkg_err.AddInvalidParam(fmt.Errorf("operator '%s' not valid", col), col,
					"operator %v not not valid", col)
				return "", err
			}

			operator = v

			switch operator {

			case "DATE":
				operator = " = "
				arr[i] = operator
				arr[i-1] = fmt.Sprintf("DATE(%v)", arr[i-1])

			case "DATE_GTE":
				operator = " >= "
				arr[i] = operator
				arr[i-1] = fmt.Sprintf("DATE(%v)", arr[i-1])

			case "DATE_LTE":
				operator = " <= "
				arr[i] = operator
				arr[i-1] = fmt.Sprintf("DATE(%v)", arr[i-1])

			}

		default:
			if sqlInjection(v) {
				err := pkg_err.AddInvalidParam(fmt.Errorf("value of column '%s' not valid", col), col,
					"value of column %v not not valid", col)
				return "", err
			}
			col = ""
			operator = ""
		}
	}

	result = ""
	for _, v := range arr {
		if column[v] != "" {
			result += " " + column[v]
		} else {
			result += " " + v
		}
	}

	result = "(" + result + ")"

	return result, nil
}

func sqlInjection(value string) (status bool) {

	arr := strings.Split(strings.ToUpper(strings.ReplaceAll(value, "'", "")), " ")
	for _, v := range arr {
		if sqlKeywords[v] {
			return true
		}
	}

	return false
}
