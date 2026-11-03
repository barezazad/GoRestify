package pkg_sql

import (
	"fmt"
	"strconv"
	"strings"

	"GoRestify/pkg/pkg_err"
)

// Parser converts the filter DSL into a parameterized WHERE fragment.
// Columns and operators are allow-listed; values become ? placeholders.
// Supported: [eq] [ne] [gt] [lt] [gte] [lte] [like] [and] [or] [date] [date_gte] [date_lte] and ( ).
func Parser(str string, cols []string) (query string, args []interface{}, err error) {
	column := buildColumnMap(cols)

	normalized := strings.ReplaceAll(str, " ", "_SPACE_")
	normalized = strings.ReplaceAll(normalized, "\"", "'")
	normalized = strings.ReplaceAll(normalized, "(", "PARENTHESES_OPEN ")
	normalized = strings.ReplaceAll(normalized, ")", " PARENTHESES_CLOSE")

	for k, v := range filterKeys {
		normalized = strings.ReplaceAll(normalized, k, v)
	}

	tokens := strings.Split(normalized, " ")
	if len(tokens) == 0 {
		return "", nil, fmt.Errorf("filter is not valid")
	}

	var parts []string
	var col, opToken string

	for _, raw := range tokens {
		token := strings.TrimSpace(raw)
		token = strings.ReplaceAll(token, "_SPACE_", " ")
		token = strings.ReplaceAll(token, "PARENTHESES_OPEN", "(")
		token = strings.ReplaceAll(token, "PARENTHESES_CLOSE", ")")

		if token == "" {
			continue
		}

		if token == "(" || token == ")" {
			parts = append(parts, token)
			continue
		}

		if logicalOps[token] && col == "" {
			parts = append(parts, token)
			continue
		}

		switch {
		case col == "":
			resolved, ok := column[token]
			if !ok {
				err = pkg_err.AddInvalidParam(fmt.Errorf("column '%s' not valid", token), token,
					"column %v not not valid", token)
				return "", nil, err
			}
			col = resolved

		case opToken == "" && col != "":
			if _, ok := comparisonOps[token]; !ok {
				err = pkg_err.AddInvalidParam(fmt.Errorf("operator '%s' not valid", token), token,
					"operator %v not not valid", token)
				return "", nil, err
			}
			opToken = token

		default:
			value := stripQuotes(token)
			if isUnsafeFilterValue(value) {
				err = pkg_err.AddInvalidParam(fmt.Errorf("value of column '%s' not valid", col), col,
					"value of column %v not not valid", col)
				err = pkg_err.SetCustom(err, pkg_err.ValidationFailedErr)
				return "", nil, err
			}

			sqlOp := comparisonOps[opToken]
			left := col

			switch opToken {

			case "DATE", "DATE_GTE", "DATE_LTE":
				left = fmt.Sprintf("DATE(%s)", col)

			case "LIKE":
				value = "%" + strings.Trim(value, "%") + "%"
			}

			parts = append(parts, left, sqlOp, "?")
			args = append(args, filterValue(opToken, col, value))
			col = ""
			opToken = ""
		}
	}

	if col != "" || opToken != "" {
		err = pkg_err.AddInvalidParam(fmt.Errorf("filter is not valid"), "filter",
			"filter is not valid")
		return "", nil, err
	}

	return "(" + strings.Join(parts, " ") + ")", args, nil
}

// filterValue binds id columns as integers so MySQL can use the numeric index.
// Other values stay strings. LIKE and date operators always stay strings.
func filterValue(opToken, col, value string) interface{} {
	if opToken == "LIKE" || opToken == "DATE" || opToken == "DATE_GTE" || opToken == "DATE_LTE" || !idColumn(col) {
		return value
	}
	if len(value) > 1 && value[0] == '0' {
		return value
	}
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return value
	}
	return n
}

func idColumn(col string) bool {
	name := col
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	name = strings.Trim(name, "`")
	return name == "id" || strings.HasSuffix(name, "_id")
}

func buildColumnMap(cols []string) map[string]string {
	column := make(map[string]string)

	for _, v := range cols {
		if strings.Contains(v, " as ") || strings.Contains(v, " AS ") {
			splitString := strings.Split(v, " ")
			if len(splitString) >= 3 {
				column[splitString[2]] = splitString[0]
			}
		}

		if strings.Contains(v, ".") {
			splitString := strings.Split(v, ".")
			column[splitString[len(splitString)-1]] = v
			column[v] = v
		} else {
			column[v] = v
		}
	}

	return column
}

func stripQuotes(value string) string {
	if len(value) >= 2 {
		if (value[0] == '\'' && value[len(value)-1] == '\'') ||
			(value[0] == '"' && value[len(value)-1] == '"') {
			return value[1 : len(value)-1]
		}
	}
	return value
}
