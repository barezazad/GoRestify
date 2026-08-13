package account_type

import "GoRestify/pkg/pkg_types"

// account type enum
const (
	User     pkg_types.Enum = "user"
	Customer pkg_types.Enum = "customer"
)

// List account type list
var List = []pkg_types.Enum{
	User,
	Customer,
}
