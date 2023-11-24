package param

import (
	"GoRestify/pkg/pkg_sql"
)

// parseFilter call parser for convert urlQuery to SQL query
func (p *Param) parseFilter(cols []string) (result string, err error) {

	if p.Filter == "" {
		return
	}

	if result, err = pkg_sql.Parser(p.Filter, cols); err != nil {
		return
	}

	return
}
