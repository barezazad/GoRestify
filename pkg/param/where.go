package param

import (
	"strings"
)

func (p *Param) parseWhere(cols []string) (whereStr string, args []interface{}, err error) {
	var whereArr []string

	var resultFilter string
	var filterArgs []interface{}
	if resultFilter, filterArgs, err = p.parseFilter(cols); err != nil {
		return
	}

	if resultFilter != "" {
		whereArr = append(whereArr, resultFilter)
		args = append(args, filterArgs...)
	}

	if !p.PreCondition.Empty() {
		if err = validateCondition(p.PreCondition.Query, p.PreCondition.Args); err != nil {
			return "", nil, err
		}
		whereArr = append(whereArr, "("+p.PreCondition.Query+")")
		args = append(args, p.PreCondition.Args...)
	}

	if len(whereArr) > 0 {
		whereStr = strings.Join(whereArr, " AND ")
	}

	return
}

// ParseWhere combines preConditions and filter into a parameterized WHERE clause.
func (p *Param) ParseWhere(cols []string) (whereStr string, args []interface{}, err error) {
	return p.parseWhere(cols)
}

// ParseWhereDelete is used when the table has deleted_at column.
// User/filter clauses are wrapped so OR filters cannot bypass soft-delete.
func (p *Param) ParseWhereDelete(cols []string) (whereStr string, args []interface{}, err error) {
	if whereStr, args, err = p.parseWhere(cols); err != nil {
		return
	}

	whereStr = strings.TrimSpace(whereStr)
	if whereStr == "" {
		whereStr = "deleted_at IS NULL"
		return
	}

	whereStr = "deleted_at IS NULL AND (" + whereStr + ")"
	return
}
