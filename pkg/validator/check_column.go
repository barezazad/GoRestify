package validator

import (
	"strings"

	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/utils"
)

// CheckColumns will check columns for security
func CheckColumns(cols []string, requestedCols string) (string, error) {
	var err error

	if requestedCols == "*" || requestedCols == "" {
		return strings.Join(cols, ","), nil
	}

	column := make(map[string]string)

	for _, v := range cols {
		if strings.Contains(v, " as ") {
			splitCol := strings.Split(v, " as ")
			column[splitCol[1]] = v
		} else if strings.Contains(v, " AS ") {
			splitCol := strings.Split(v, " AS ")
			column[splitCol[1]] = v
		} else {
			splitCol := strings.Split(v, ".")
			column[splitCol[1]] = v
		}
	}

	variates := strings.Split(requestedCols, ",")
	requestedCols = ""
	for i, v := range variates {
		if column[v] != "" {
			v = column[v]
		}
		if i == 0 {
			requestedCols += v
		} else {
			requestedCols += "," + v
		}

		if ok, _ := utils.ArrayIncludes(cols, v); !ok {
			err = pkg_err.AddInvalidParam(err, v, pkg_err.XisNotValid, v)
			err = pkg_err.SetCustom(err, pkg_err.ValidationFailedErr)
		}
	}

	return requestedCols, err

}
