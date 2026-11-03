package param

import (
	"fmt"
	"regexp"
	"strings"
)

// matchOperators is longest-first so date_gte wins over date and gte wins over gt.
var matchOperators = []string{
	"date_gte",
	"date_lte",
	"like",
	"date",
	"gte",
	"lte",
	"eq",
	"ne",
	"gt",
	"lt",
}

var matchColumnRE = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(\.[A-Za-z_][A-Za-z0-9_]*)*$`)

// translateMatch rewrites the WAF-safe match DSL into the internal bracket DSL.
// A clause is column.operator.value. A comma is AND and ~or~ is OR.
// status.eq.active,city_id.eq.12 → status[eq]'active'[and]city_id[eq]'12'
func translateMatch(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}

	clauses, joins, err := splitMatch(raw)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	for i, clause := range clauses {
		dsl, err := matchClause(clause)
		if err != nil {
			return "", err
		}
		if i > 0 {
			b.WriteString(joins[i-1])
		}
		b.WriteString(dsl)
	}
	return b.String(), nil
}

func splitMatch(raw string) (clauses, joins []string, err error) {
	start := 0
	for i := 0; i < len(raw); {
		switch {
		case strings.HasPrefix(raw[i:], "~or~"):
			clauses = append(clauses, raw[start:i])
			joins = append(joins, "[or]")
			i += len("~or~")
			start = i
		case raw[i] == ',':
			clauses = append(clauses, raw[start:i])
			joins = append(joins, "[and]")
			i++
			start = i
		default:
			i++
		}
	}
	clauses = append(clauses, raw[start:])

	if len(clauses) == 0 {
		return nil, nil, fmt.Errorf("match is not valid")
	}
	return clauses, joins, nil
}

func matchClause(clause string) (string, error) {
	clause = strings.TrimSpace(clause)
	op, dot := findMatchOperator(clause)
	if op == "" {
		return "", fmt.Errorf("match is not valid")
	}

	column := clause[:dot]
	value := clause[dot+len(op)+2:]
	if !matchColumnRE.MatchString(column) || value == "" {
		return "", fmt.Errorf("match is not valid")
	}
	if strings.ContainsAny(value, `'"`+"[]") {
		return "", fmt.Errorf("match is not valid")
	}

	return fmt.Sprintf("%s[%s]'%s'", column, op, value), nil
}

func findMatchOperator(clause string) (op string, dot int) {
	dot = -1
	for _, candidate := range matchOperators {
		needle := "." + candidate + "."
		idx := strings.Index(clause, needle)
		if idx <= 0 {
			continue
		}
		if dot == -1 || idx < dot {
			dot = idx
			op = candidate
		}
	}
	return op, dot
}
