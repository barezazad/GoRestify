package account_status

import "GoRestify/pkg/pkg_types"

// account status enum
const (
	Active   pkg_types.Enum = "active"
	InActive pkg_types.Enum = "in_active"
)

// List account status list
var List = []pkg_types.Enum{
	Active,
	InActive,
}
