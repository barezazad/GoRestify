package pkg_sql

import (
	"regexp"
	"strings"
)

// filterKeys maps DSL tokens to internal markers.
var (
	filterKeys = map[string]string{
		"[eq]":       " = ",
		"[ne]":       " != ",
		"[gt]":       " > ",
		"[lt]":       " < ",
		"[gte]":      " >= ",
		"[lte]":      " <= ",
		"[like]":     " LIKE ",
		"[and]":      " AND ",
		"[or]":       " OR ",
		"[date]":     " DATE ",
		"[date_gte]": " DATE_GTE ",
		"[date_lte]": " DATE_LTE ",
	}

	// comparisonOps: allow-listed value-binding operators → SQL operator text.
	comparisonOps = map[string]string{
		"=":        "=",
		"!=":       "!=",
		">":        ">",
		"<":        "<",
		">=":       ">=",
		"<=":       "<=",
		"LIKE":     "LIKE",
		"DATE":     "=",
		"DATE_GTE": ">=",
		"DATE_LTE": "<=",
	}

	logicalOps = map[string]bool{
		"AND": true,
		"OR":  true,
	}

	// unsafeFilterValue matches injection-like payloads that must be rejected (never sent to SQL).
	unsafeFilterValue = regexp.MustCompile(`(?i)` + strings.Join([]string{
		`;`,   // stacked queries
		`--`,  // SQL comment
		`/\*`, // block comment start
		`\*/`, // block comment end
		`'`,   // quote breakout
		`"`,   // quote breakout
		`\bSELECT\b`,
		`\bINSERT\b`,
		`\bUPDATE\b`,
		`\bDELETE\b`,
		`\bDROP\b`,
		`\bUNION\b`,
		`\bEXEC\b`,
		`\bEXECUTE\b`,
		`\bTRUNCATE\b`,
		`\bALTER\b`,
		`\bCREATE\b`,
		`\bSLEEP\s*\(`,
		`\bBENCHMARK\s*\(`,
		`\bOR\b`,
		`\bAND\b`,
	}, "|"))
)

// isUnsafeFilterValue reports whether a filter value looks like an SQL injection attempt.
func isUnsafeFilterValue(value string) bool {
	return unsafeFilterValue.MatchString(value)
}
