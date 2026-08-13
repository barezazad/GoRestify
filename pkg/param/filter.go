package param

import (
	"GoRestify/pkg/pkg_sql"
)

// parseFilter converts the URL filter DSL into a parameterized SQL fragment.
func (p *Param) parseFilter(cols []string) (result string, args []interface{}, err error) {
	if p.Filter == "" {
		return
	}

	return pkg_sql.Parser(p.Filter, cols)
}
