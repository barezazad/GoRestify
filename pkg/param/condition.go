package param

import (
	"fmt"
	"strings"

	"GoRestify/pkg/pkg_err"
)

// Condition is a parameterized SQL fragment (query uses ? placeholders).
type Condition struct {
	Query string
	Args  []interface{}
}

// Empty reports whether the condition has no query.
func (c Condition) Empty() bool {
	return strings.TrimSpace(c.Query) == ""
}

// SetPreCondition sets a developer-owned parameterized precondition.
// query must use ? for values; args count must match; unsafe SQL tokens are rejected.
func (p *Param) SetPreCondition(query string, args ...interface{}) error {
	query = strings.TrimSpace(query)
	if err := validateCondition(query, args); err != nil {
		return err
	}
	p.PreCondition = Condition{Query: query, Args: args}
	return nil
}

// AddPreCondition ANDs another parameterized condition onto PreCondition.
func (p *Param) AddPreCondition(query string, args ...interface{}) error {
	query = strings.TrimSpace(query)
	if err := validateCondition(query, args); err != nil {
		return err
	}

	if p.PreCondition.Empty() {
		p.PreCondition = Condition{Query: query, Args: args}
		return nil
	}

	p.PreCondition.Query = "(" + p.PreCondition.Query + ") AND (" + query + ")"
	p.PreCondition.Args = append(p.PreCondition.Args, args...)
	return nil
}

func validateCondition(query string, args []interface{}) error {
	if query == "" {
		err := pkg_err.AddInvalidParam(fmt.Errorf("precondition is empty"), "pre_condition",
			"precondition is empty")
		return pkg_err.SetCustom(err, pkg_err.ValidationFailedErr)
	}

	if strings.Contains(query, ";") ||
		strings.Contains(query, "--") ||
		strings.Contains(query, "/*") ||
		strings.Contains(query, "*/") {
		err := pkg_err.AddInvalidParam(fmt.Errorf("precondition contains unsafe SQL"), "pre_condition",
			"precondition contains unsafe SQL")
		return pkg_err.SetCustom(err, pkg_err.ValidationFailedErr)
	}

	upper := strings.ToUpper(query)
	for _, banned := range []string{
		" UNION ", " INSERT ", " UPDATE ", " DELETE ", " DROP ",
		" ALTER ", " CREATE ", " TRUNCATE ", " EXEC ", " EXECUTE ",
		" INTO OUTFILE", " LOAD_FILE", " INFORMATION_SCHEMA",
	} {
		if strings.Contains(upper, banned) {
			err := pkg_err.AddInvalidParam(fmt.Errorf("precondition contains unsafe SQL"), "pre_condition",
				"precondition contains unsafe SQL")
			return pkg_err.SetCustom(err, pkg_err.ValidationFailedErr)
		}
	}

	placeholderCount := countPlaceholders(query)
	if placeholderCount != len(args) {
		err := pkg_err.AddInvalidParam(
			fmt.Errorf("precondition placeholder count mismatch: got %d args for %d ?", len(args), placeholderCount),
			"pre_condition",
			"precondition placeholder count mismatch",
		)
		return pkg_err.SetCustom(err, pkg_err.ValidationFailedErr)
	}

	return nil
}

func countPlaceholders(query string) int {
	count := 0
	inSingle := false
	inDouble := false

	for i := 0; i < len(query); i++ {
		ch := query[i]
		switch {
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
		case ch == '"' && !inSingle:
			inDouble = !inDouble
		case ch == '?' && !inSingle && !inDouble:
			count++
		}
	}

	return count
}
